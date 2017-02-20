package sat

// trail is the state of the solver that contains the list of literals
// and their current value.
type trail []trailElem

type trailElem struct {
	Lit      Literal
	Decision bool
}

func (t *trail) Assert(l Literal, d bool) {
	*t = append(*t, trailElem{
		Lit:      l,
		Decision: d,
	})
}
