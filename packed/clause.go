package packed

import (
	"fmt"
	"sort"
	"strings"
)

// Clause represents a clause in a formula.
type Clause struct {
	lits   []Lit // lits is the list of literals
	maxVar int   // maxVar is the largest variable in Lits
}

// NewClause creates a new clause with pre-allocated space for cap literals.
// If you don't know how many literals, you can leave cap as 0 and it will
// grow automatically.
func NewClause(cap int) *Clause {
	var lits []Lit
	if cap > 0 {
		lits = make([]Lit, 0, cap)
	}

	return &Clause{
		lits: lits,
	}
}

// Lits are the literals in this clause.
func (c *Clause) Lits() []Lit {
	return c.lits
}

// MaxVar returns the maximum value variable within this clause.
func (c *Clause) MaxVar() int {
	return c.maxVar
}

// SetLits directly sets the lits slice as the lits within this Clause.
// It is unsafe to use lits for anything else after this.
func (c *Clause) SetLits(lits []Lit) {
	c.lits = lits

	// Sort
	sort.Slice(lits, func(i, j int) bool {
		return lits[i] < lits[j]
	})

	// Update the max var
	for _, l := range c.lits {
		if v := l.Var(); v > c.maxVar {
			c.maxVar = v
		}
	}
}

// Ref returns a reference to this clause.
func (c *Clause) Ref() *Clause {
	return c
}

// String representation
func (c *Clause) String() string {
	ls := make([]string, len(c.lits))
	for i, l := range c.lits {
		ls[i] = l.String()
	}

	return fmt.Sprintf("[%s]", strings.Join(ls, ", "))
}
