package sat

import (
	"github.com/mitchellh/go-sat/cnf"
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

func (s *Solver) assertLiteral(l cnf.Literal, d bool) {
	// TODO: old legacy
	s.m.Assert(l, d)

	// If this is a decision literal, then create a new decision level
	if d {
		s.newDecisionLevel()
	}

	// Store the literal in the trail
	pl := l.Pack()
	s.assigns[pl.Var()] = BoolToTri(!pl.Sign())
	s.trail = append(s.trail, pl)
}

// newDecisionLevel creates a new decision level within the trail
func (s *Solver) newDecisionLevel() {
	s.trailIdx = append(s.trailIdx, len(s.trail))
}

// decisionLevel returns the current decision level
func (s *Solver) decisionLevel() int {
	return len(s.trailIdx)
}
