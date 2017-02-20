package sat

import (
	"fmt"
	"strings"

	"github.com/mitchellh/go-sat/cnf"
)

// trail is the state of the solver that contains the list of literals
// and their current value.
type trail []trailElem

type trailElem struct {
	Lit      cnf.Literal
	Decision bool
}

// DecisionsLen returns the number of decision variables are in the trail.
func (t trail) DecisionsLen() int {
	count := 0
	for _, e := range t {
		if e.Decision {
			count++
		}
	}

	return count
}

// TrimToLastDecision trims the trail to the last decision (but not including
// it) and returns the last decision literal.
func (t *trail) TrimToLastDecision() cnf.Literal {
	var i int
	for i = len(*t) - 1; i >= 0; i-- {
		if (*t)[i].Decision {
			break
		}
	}

	result := (*t)[i].Lit
	*t = (*t)[:i]
	return result
}

// String returns human readable output for a trail that shows the
// literals chosen. Decision literals are prefixed with '|'.
func (t trail) String() string {
	result := make([]string, len(t))
	for i, e := range t {
		v := ""
		if e.Decision {
			v = "| "
		}

		v += fmt.Sprintf("%d", e.Lit)
		result[i] = v
	}

	return "[" + strings.Join(result, ", ") + "]"
}

// Assert adds the new literal to the trail.
func (t *trail) Assert(l cnf.Literal, d bool) {
	*t = append(*t, trailElem{
		Lit:      l,
		Decision: d,
	})
}

// IsFormulaFalse returns true if the given Formula f is false in the
// current valuation (trail).
func (t trail) IsFormulaFalse(f cnf.Formula) bool {
	// If we have no trail, we can't contain the negated formula
	if len(t) == 0 {
		return false
	}

	// We need to find ONE negated clause in f
	for _, c := range f {
		if t.IsClauseFalse(c) {
			return true
		}
	}

	return false
}

func (t trail) IsClauseFalse(c cnf.Clause) bool {
	for _, l := range c {
		if !t.IsLiteralFalse(l) {
			return false
		}
	}

	return true
}

func (t trail) IsLiteralFalse(l cnf.Literal) bool {
	l = l.Negate()
	for _, e := range t {
		if e.Lit == l {
			return true
		}
	}

	return false
}

func (t trail) IsLiteralTrue(l cnf.Literal) bool {
	for _, e := range t {
		if e.Lit == l {
			return true
		}
	}

	return false
}
