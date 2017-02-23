package sat

import (
	"fmt"
	"strings"

	"github.com/mitchellh/go-sat/cnf"
)

// trail is the state of the solver that contains the list of literals
// and their current value.
type trail struct {
	elems       []trailElem
	set         map[cnf.Literal]int
	decisionLen int
}

type trailElem struct {
	Lit      cnf.Literal
	Decision bool
}

func newTrail(cap int) *trail {
	return &trail{
		elems: make([]trailElem, 0, cap),
		set:   make(map[cnf.Literal]int),
	}
}

// Len returns the number of variables are in the trail.
func (t *trail) Len() int {
	return len(t.elems)
}

func (t *trail) CurrentLevel() int {
	return t.decisionLen
}

// DecisionsLen returns the number of decision variables are in the trail.
func (t *trail) DecisionsLen() int {
	return t.decisionLen
}

// Assert adds the new literal to the trail.
func (t *trail) Assert(l cnf.Literal, d bool) {
	// Add it to the list
	t.elems = append(t.elems, trailElem{
		Lit:      l,
		Decision: d,
	})

	// If this is a decision var, then incr our counter cache
	if d {
		t.decisionLen++
	}

	// Store it in our set
	t.set[l] = t.decisionLen
}

// TrimToLastDecision trims the trail to the last decision (but not including
// it) and returns the last decision literal.
func (t *trail) TrimToLastDecision() cnf.Literal {
	var i int
	for i = len(t.elems) - 1; i >= 0; i-- {
		if t.elems[i].Decision {
			break
		}
	}

	for _, e := range t.elems[i:] {
		delete(t.set, e.Lit)
	}

	result := t.elems[i].Lit
	t.elems = t.elems[:i]
	t.decisionLen--
	return result
}

func (t *trail) TrimToLevel(level int) {
	var i int
	for i = len(t.elems) - 1; i >= 0; i-- {
		if t.set[t.elems[i].Lit] <= level {
			break
		}
	}

	for _, e := range t.elems[i+1:] {
		delete(t.set, e.Lit)
	}

	t.elems = t.elems[:i+1]
	t.decisionLen = level
}

// LastAssertedLiteral returns the last asserted literal (decision or not)
// from the clause c. "Last" is defined as most recent in the trail.
func (t *trail) LastAssertedLiteral(c cnf.Clause) cnf.Literal {
	for i := len(t.elems) - 1; i >= 0; i-- {
		lit := t.elems[i].Lit
		for _, l := range c {
			if lit == l {
				return lit
			}
		}
	}

	return cnf.Literal(0)
}

// Level returns the decision level that the given literal exists at.
func (t *trail) Level(l cnf.Literal) int {
	return t.set[l]
}

func (t *trail) MaxLevel(c cnf.Clause) int {
	max := 0
	for _, l := range c {
		if v := t.set[l]; v > max {
			max = v
		}
	}

	return max
}

// String returns human readable output for a trail that shows the
// literals chosen. Decision literals are prefixed with '|'.
func (t trail) String() string {
	result := make([]string, len(t.elems))
	for i, e := range t.elems {
		v := ""
		if e.Decision {
			v = "| "
		}

		v += fmt.Sprintf("%d", e.Lit)
		result[i] = v
	}

	return "[" + strings.Join(result, ", ") + "]"
}

// IsUnit returns true if the clause c is a unit clause in t with
// literal l. Clause c must be a clause within the formula that this
// trail is being used for.
func (t *trail) IsUnit(c cnf.Clause, unitL cnf.Literal) bool {
	m := t.set

	// If we already have the unit literal we're looking for (+ or -),
	// then this is not a unit clause
	if _, ok := m[unitL]; ok {
		return false
	}
	if _, ok := m[unitL.Negate()]; ok {
		return false
	}

	for _, l := range c {
		if l == unitL || l == unitL.Negate() {
			continue
		}

		if _, ok := m[l.Negate()]; !ok {
			return false
		}
	}

	return true
}

// IsFormulaFalse returns a non-zero Clause if the given Formula f is
// false in the current valuation (trail). This non-zero clause is a false
// clause.
func (t *trail) IsFormulaFalse(f cnf.Formula) cnf.Clause {
	// If we have no trail, we can't contain the negated formula
	if len(t.elems) == 0 {
		return cnf.Clause(nil)
	}

	// We need to find ONE negated clause in f
	for _, c := range f {
		if t.IsClauseFalse(c) {
			return c
		}
	}

	return cnf.Clause(nil)
}

func (t *trail) IsClauseFalse(c cnf.Clause) bool {
	for _, l := range c {
		if !t.IsLiteralFalse(l) {
			return false
		}
	}

	return true
}

func (t *trail) IsLiteralFalse(l cnf.Literal) bool {
	l = l.Negate()
	for _, e := range t.elems {
		if e.Lit == l {
			return true
		}
	}

	return false
}

func (t *trail) IsLiteralTrue(l cnf.Literal) bool {
	for _, e := range t.elems {
		if e.Lit == l {
			return true
		}
	}

	return false
}
