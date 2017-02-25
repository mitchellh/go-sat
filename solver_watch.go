package sat

import (
	"github.com/mitchellh/go-sat/packed"
)

type watcher struct {
	Clause *packed.Clause
	Lit    packed.Lit
}

func (s *Solver) watchClause(c *packed.Clause) {
	lits := c.Lits()
	c0 := lits[0].Neg()
	c1 := lits[1].Neg()

	s.watches[c0] = append(s.watches[c0], &watcher{
		Clause: c,
		Lit:    lits[1],
	})
	s.watches[c1] = append(s.watches[c1], &watcher{
		Clause: c,
		Lit:    lits[0],
	})
}

func (s *Solver) propagate() *packed.Clause {
	var conflict *packed.Clause

	for s.qhead < len(s.trail) {
		p := s.trail[s.qhead]
		s.qhead++

		watches := s.watches[p]
		i := 0
		j := 0
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

			// We keep the false literal in w[1]
			iLits := iW.Clause.Lits()
			if iLits[0] == p.Neg() {
				iLits[0], iLits[1] = iLits[1], p.Neg()
			}

			// If first watch is true, clause is satisfied
			newW := &watcher{Clause: iW.Clause, Lit: iLits[0]}
			if iLits[0] != iW.Lit && s.ValueLit(iLits[0]) == True {
				watches[j] = newW
				j++
				continue
			}

			// Look for a new watch
			for k := 2; k < len(iLits); k++ {
				if s.ValueLit(iLits[k]) != False {
					iLits[1], iLits[k] = iLits[k], p.Neg()
					i1 := iLits[1].Neg()
					s.watches[i1] = append(s.watches[i1], newW)
					goto PROP_LOOP
				}
			}

			// No watch found, clause is unit
			watches[j] = newW
			j++
			if s.ValueLit(iLits[0]) == False {
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
					iLits[0], iW.Clause)
			}

			// Enqueue
			s.assertLiteral(iLits[0], iW.Clause)

		PROP_LOOP:
		}

		s.watches[p] = watches[:j]
	}

	return conflict
}
