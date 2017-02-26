package sat

import (
	"github.com/mitchellh/go-sat/cnf"
)

// Solver is a SAT solver. This should be created with New to get
// the proper internal memory allocations. Using a manually allocated
// Solver will probably crash.
type Solver struct {
	// Trace, if set to true, will output trace debugging information
	// via the standard library `log` package. If true, Tracer must also
	// be set to a non-nil value. The easiest implmentation is a logger
	// created with log.NewLogger.
	Trace  bool
	Tracer Tracer

	// decideLiterals is to be set by tests to force a certain decision
	// literal ordering. This can be used to exercise specific solver
	// behavior being tested.
	decideLiterals []int

	//---------------------------------------------------------------
	// Internal fields, do not set
	//---------------------------------------------------------------
	result satResult

	// problem
	clauses []cnf.Clause     // clauses to solve
	vars    map[int]struct{} // list of available vars

	// two-literal watching
	qhead   int
	watches map[cnf.Lit][]*watcher
	seen    map[int]int8

	//---------------------------------------------------------------
	// trail
	//---------------------------------------------------------------

	// trail is the actual trail of assigned literals. The value assigned
	// is in the assigns map.
	trail []cnf.Lit

	// trailIdx keeps track of the indices for new decision levels.
	// trailIdx[level] = index to the start of that level in trail
	trailIdx []int

	// assigns keeps track of variable assignment values. unassigned variables
	// are never present in assigns.
	assigns map[int]Tribool

	// varinfo holds information about an assigned variable. unassigned
	// variables may be present here but their resulting information is
	// garbage.
	varinfo map[int]varinfo
}

// New creates a new solver and allocates the basics for it.
func New() *Solver {
	return &Solver{
		result: satResultUndef,

		// problem
		vars: make(map[int]struct{}),

		// trail
		assigns: make(map[int]Tribool),
		varinfo: make(map[int]varinfo),

		// two-literal watches
		watches: make(map[cnf.Lit][]*watcher),
		seen:    make(map[int]int8),
	}
}

// Solve finds a solution for the formula, returning true on satisfiability.
func (s *Solver) Solve() bool {
	if s.Trace {
		s.Tracer.Printf("[TRACE] sat: starting solve()")
	}

	// Check the result. This can be set already by a prior call to Solve
	// or via the AddClause process.
	if s.result != satResultUndef {
		if s.Trace {
			s.Tracer.Printf(
				"[TRACE] sat: result is already available: %s", s.result)
		}

		return s.result == satResultSat
	}

	for {
		if s.Trace {
			s.Tracer.Printf("[TRACE] sat: new iteration. trail: %s", s.trailString())
		}

		conflictC := s.propagate()
		if conflictC != nil {
			if s.Trace {
				s.Tracer.Printf("[TRACE] sat: current trail contains negated formula. trail: %s", s.trailString())
				s.Tracer.Printf("[TRACE] sat: conflict clause: %s", conflictC)
			}

			// If we have no more decisions within the trail, then we've
			// failed finding a satisfying value.
			if s.decisionLevel() == 0 {
				if s.Trace {
					s.Tracer.Printf("[TRACE] sat: at decision level 0. UNSAT")
				}

				return false
			}

			// Learn
			learnt, level := s.learn(conflictC)
			if s.Trace {
				s.Tracer.Printf("[TRACE] sat: learned clause: %s", learnt)
			}

			// Backjump
			s.trimToDecisionLevel(level)
			if s.Trace {
				s.Tracer.Printf(
					"[TRACE] sat: backjump to level %d, trail: %s",
					level, s.trail)
			}

			// Add our learned clause
			lits := learnt
			if s.Trace {
				s.Tracer.Printf(
					"[TRACE] sat: asserting learned literal: %s", lits[0])
			}
			if len(lits) == 1 {
				s.assertLiteral(lits[0], nil)
			} else {
				c := cnf.Clause(lits)
				s.clauses = append(s.clauses, c)
				s.watchClause(c)
				s.assertLiteral(lits[0], c)
			}
		} else {
			// Choose a literal to assert.
			lit := s.selectLiteral()

			// If it is undef it means there are no more literals which means
			// we have solved the formula
			if lit == cnf.LitUndef {
				if s.Trace {
					s.Tracer.Printf("[TRACE] sat: solver found solution: %s", s.trail)
				}

				return true
			}

			if s.Trace {
				s.Tracer.Printf("[TRACE] sat: assert: %s (decision)", lit)
			}
			s.newDecisionLevel()
			s.assertLiteral(lit, nil)
		}
	}

	return false
}

func (s *Solver) selectLiteral() cnf.Lit {
	for raw, _ := range s.vars {
		if _, ok := s.assigns[raw]; !ok {
			return cnf.NewLit(raw, false)
		}
	}

	return cnf.LitUndef
}

//-------------------------------------------------------------------
// Private types
//-------------------------------------------------------------------

type satResult byte

const (
	satResultUndef satResult = iota
	satResultUnsat
	satResultSat
)

type varinfo struct {
	reason cnf.Clause
	level  int
}

// Tribool is a tri-state boolean with undefined as the 3rd state.
type Tribool uint8

const (
	True  Tribool = 0
	False         = 1
	Undef         = 2
)

func BoolToTri(b bool) Tribool {
	if b {
		return True
	}

	return False
}
