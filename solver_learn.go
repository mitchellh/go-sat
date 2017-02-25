package sat

import (
	"github.com/mitchellh/go-sat/packed"
)

func (s *Solver) learn(c packed.Clause) (packed.Clause, int) {
	// Determine our learned clause
	pathC := 0
	p := packed.LitUndef
	learnt := make([]packed.Lit, 1)
	idx := len(s.trail) - 1
	for {
		j := 0
		if p != packed.LitUndef {
			j = 1
		}

		for ; j < len(c); j++ {
			q := c[j]
			qVar := q.Var()
			qLevel := s.level(qVar)
			if s.seen[qVar] == 0 && qLevel > 0 {
				s.seen[qVar] = 1
				if qLevel >= s.decisionLevel() {
					pathC++
				} else {
					learnt = append(learnt, q)
				}
			}
		}

		// Select next clause
		for s.seen[s.trail[idx].Var()] == 0 {
			idx--
		}
		idx--

		p = s.trail[idx+1]
		c = s.varinfo[p.Var()].reason
		s.seen[p.Var()] = 0
		pathC--

		if pathC <= 0 {
			break
		}
	}
	learnt[0] = p.Neg()

	// Determine the backtrack level
	backjumpLevel := 0
	if len(learnt) > 1 {
		maxI := 1
		maxLevel := s.level(learnt[maxI].Var())
		for i := 2; i < len(learnt); i++ {
			if l := s.level(learnt[i].Var()); l > maxLevel {
				maxI = i
				maxLevel = l
			}
		}

		learnt[maxI], learnt[1] = learnt[1], learnt[maxI]
		backjumpLevel = maxLevel
	}

	// Clear seen for learnt clause
	for _, l := range learnt {
		s.seen[l.Var()] = 0
	}

	return packed.Clause(learnt), backjumpLevel
}
