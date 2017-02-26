package sat

import (
	"github.com/mitchellh/go-sat/cnf"
)

// learn performs the clause learning process after a conflict is found.
// This returns the learned clause as well as the level to backjump to.
func (s *Solver) learn(c cnf.Clause) int {
	// Determine our learned clause
	pathC := 0
	s.learned = s.learned[:1]
	p := cnf.LitUndef
	idx := len(s.trail) - 1
	for {
		j := 0
		if p != cnf.LitUndef {
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
					s.learned = append(s.learned, q)
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
	s.learned[0] = p.Neg()

	// Determine the level to backjump to. This is simply the maximum
	// level represented in our learned clause.
	backjumpLevel := 0
	if len(s.learned) > 1 {
		maxI := 1
		maxLevel := s.level(s.learned[maxI].Var())
		for i := 2; i < len(s.learned); i++ {
			if l := s.level(s.learned[i].Var()); l > maxLevel {
				maxI = i
				maxLevel = l
			}
		}

		s.learned[maxI], s.learned[1] = s.learned[1], s.learned[maxI]
		backjumpLevel = maxLevel
	}

	// Clear seen for learned clause so that learning can visit them
	// again on the next go-around.
	for _, l := range s.learned {
		s.seen[l.Var()] = 0
	}

	return backjumpLevel
}
