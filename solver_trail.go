package sat

import (
	"github.com/mitchellh/go-sat/packed"
)

// This file contains the trail-related functions for the solver.

// ValueLit reads the currently set value for a literal.
func (s *Solver) ValueLit(l packed.Lit) Tribool {
	result, ok := s.assigns[l.Var()]
	if !ok || result == Undef {
		return Undef
	}

	// If the literal is negative (signed), then XOR 1 will cause the bool
	// to flip. If result is undef, this has no affect.
	if l.Sign() {
		result ^= 1
	}

	return result
}

func (s *Solver) assertLiteral(l packed.Lit) {
	// Store the literal in the trail
	v := l.Var()
	s.assigns[v] = BoolToTri(!l.Sign())
	s.varinfo[v] = varinfo{level: s.decisionLevel()}
	s.trail = append(s.trail, l)
}

// level returns the level for the variable specified by v. This variable
// must be assigned for this to be correct.
func (s *Solver) level(v int) int {
	return s.varinfo[v].level
}

// IsUnit returns true if the clause c is a unit clause in t with
// literal l. Clause c must be a clause within the formula that this
// trail is being used for.
func (s *Solver) isUnit(c packed.Clause, unitL packed.Lit) bool {
	// If we already have the unit literal we're looking for (+ or -),
	// then this is not a unit clause
	if _, ok := s.assigns[unitL.Var()]; ok {
		return false
	}

	for _, l := range c.Lits() {
		if l.Var() == unitL.Var() {
			continue
		}

		if v := s.ValueLit(l); v == Undef || v == True {
			return false
		}
	}

	return true
}

// IsFormulaFalse returns a non-zero Clause if the given Formula f is
// false in the current valuation (trail). This non-zero clause is a false
// clause.
func (s *Solver) isFormulaFalse() *packed.Clause {
	// If we have no trail, we can't contain the negated formula
	if len(s.trail) == 0 {
		return nil
	}

	// We need to find ONE negated clause in f
	for _, c := range s.clauses {
		found := false
		for _, l := range c.Lits() {
			if s.ValueLit(l) != False {
				found = true
				break
			}
		}

		if !found {
			return c.Ref()
		}
	}

	return nil
}

// newDecisionLevel creates a new decision level within the trail
func (s *Solver) newDecisionLevel() {
	s.trailIdx = append(s.trailIdx, len(s.trail))
}

// decisionLevel returns the current decision level
func (s *Solver) decisionLevel() int {
	return len(s.trailIdx)
}

// trimToDecisionLevel trims the trail down to the given level (including
// that level).
func (s *Solver) trimToDecisionLevel(level int) {
	if s.decisionLevel() <= level {
		return
	}

	lastIdx := s.trailIdx[level]

	// Unassign anything in the trail in higher levels
	for i := len(s.trail) - 1; i >= lastIdx; i-- {
		delete(s.assigns, s.trail[i].Var())
	}

	// Reset the trail length
	s.trail = s.trail[:lastIdx]
	s.trailIdx = s.trailIdx[:level]
}
