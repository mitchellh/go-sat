package sat

import (
	"github.com/mitchellh/go-sat/packed"
)

type watcher struct {
	Clause packed.Clause
	Lit    packed.Lit
}

func (s *Solver) watchClause(c packed.Clause) {
	c0 := c[0].Neg()
	c1 := c[1].Neg()

	s.watches[c0] = append(s.watches[c0], &watcher{
		Clause: c,
		Lit:    c[1],
	})
	s.watches[c1] = append(s.watches[c1], &watcher{
		Clause: c,
		Lit:    c[0],
	})
}

func (s *Solver) propagate() packed.Clause {
	var conflict packed.Clause

	for s.qhead < len(s.trail) {
		p := s.trail[s.qhead]
		s.qhead++

		watches := s.watches[p]
		i := 0
		j := 0

	PROP_LOOP:
		for i < len(watches) {
			iW := watches[i]
			if s.ValueLit(iW.Lit) == True {
				watches[j] = watches[i]
				j++
				i++
				continue
			}

			// Increment
			i++

			// We're going to need ~p for the remainder
			pNeg := p.Neg()

			// We keep the false literal in w[1]
			iLits := iW.Clause
			first := iLits[0]
			if first == pNeg {
				iLits[0], iLits[1] = iLits[1], pNeg
				first = iLits[0]
			}

			// If first watch is true, clause is satisfied
			newW := &watcher{Clause: iW.Clause, Lit: first}
			if first != iW.Lit && s.ValueLit(first) == True {
				watches[j] = newW
				j++
				continue
			}

			// Look for a new watch
			for k := 2; k < len(iLits); k++ {
				if s.ValueLit(iLits[k]) != False {
					iLits[1], iLits[k] = iLits[k], pNeg
					i1 := iLits[1].Neg()
					s.watches[i1] = append(s.watches[i1], newW)
					continue PROP_LOOP
				}
			}

			// No watch found, clause is unit
			watches[j] = newW
			j++
			if s.ValueLit(first) == False {
				conflict = iW.Clause
				s.qhead = len(s.trail)

				for i < len(watches) {
					watches[j] = watches[i]
					i++
					j++
				}

				continue
			}

			if s.Trace {
				s.Tracer.Printf(
					"[TRACE] sat: asserting unit literal %s in clause %s",
					first, iW.Clause)
			}

			// Enqueue
			s.assertLiteral(first, iW.Clause)
		}

		s.watches[p] = watches[:j]
	}

	return conflict
}
