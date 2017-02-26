package sat

import (
	"fmt"
	"strings"

	"github.com/mitchellh/go-sat/cnf"
)

// This file contains the trail-related functions for the solver.

// Assignments returns the assigned variables and their value (true or false).
// This is only valid if Solve returned true, in which case this is the
// solution.
func (s *Solver) Assignments() map[int]bool {
	result := make(map[int]bool)
	for k, v := range s.assigns {
		if v != triUndef {
			result[k] = v == triTrue
		}
	}

	return result
}

// ValueLit reads the currently set value for a literal.
func (s *Solver) valueLit(l cnf.Lit) tribool {
	result, ok := s.assigns[l.Var()]
	if !ok || result == triUndef {
		return triUndef
	}

	// If the literal is negative (signed), then XOR 1 will cause the bool
	// to flip. If result is undef, this has no affect.
	if l.Sign() {
		result ^= 1
	}

	return result
}

func (s *Solver) assertLiteral(l cnf.Lit, from cnf.Clause) {
	// Store the literal in the trail
	v := l.Var()
	s.assigns[v] = boolToTri(!l.Sign())
	s.varinfo[v] = varinfo{reason: from, level: s.decisionLevel()}
	s.trail = append(s.trail, l)
}

// level returns the level for the variable specified by v. This variable
// must be assigned for this to be correct.
func (s *Solver) level(v int) int {
	return s.varinfo[v].level
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

	// Update our queue head
	s.qhead = lastIdx

	// Reset the trail length
	s.trail = s.trail[:lastIdx]
	s.trailIdx = s.trailIdx[:level]
}

// trailString is used for debugging
func (s *Solver) trailString() string {
	vs := make([]string, len(s.trail))
	for i, l := range s.trail {
		decision := ""
		for _, idx := range s.trailIdx {
			if idx == i {
				decision = "| "
				break
			}
		}

		vs[i] = fmt.Sprintf("%s%s", decision, l)
	}

	return fmt.Sprintf("[%s]", strings.Join(vs, ", "))
}
