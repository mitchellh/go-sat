package cnf

import (
	"fmt"
)

// Lit represents a literal in a formula.
//
// The least significant bit is the state of the literal
// (0 = positive, 1 = negative). The actual value of the literal is
// the literal right shifted by 1.
//
// Example: for a literal 12, value is  11000  (24)
// Example: for a literal -12, value is 11001 (25)
//
// For sorting: Lit should be sorted naturally as an integer. This will
// ensure that L and -L are always next to each other (with L coming before).
type Lit int

// LitUndef is the undefined literal. This is useful for various operations.
const LitUndef = Lit(-1)

// NewLit creates a new literal for the variable v. v must be 0..N and
// s should be true if the variable is negative.
func NewLit(v int, s bool) Lit { return Lit(v + v + boolToInt(s)) }

// NewLitInt creates a new Lit from an integer where a negative value
// implies a negated literal. Ex. -12 is the variable 12 negated.
func NewLitInt(v int) Lit {
	s := v < 0
	if s {
		v *= -1
	}

	return NewLit(v, s)
}

// Sign reads the sign of the literal. This returns true if the literal is negative.
func (l Lit) Sign() bool { return l&1 == 1 }

// Var returns the 0..N variable value for this lit.
func (l Lit) Var() int { return int(l >> 1) }

// Neg negates the literal.
func (l Lit) Neg() Lit { return Lit(l ^ 1) }

// Int returns an integer representation of this literal. +X is a true
// literal and -X is a negative literal.
func (l Lit) Int() int {
	result := l.Var()
	if l.Sign() {
		result *= -1
	}

	return result
}

// Stringer impl.
func (l Lit) String() string { return fmt.Sprintf("%d", l.Int()) }

func boolToInt(b bool) int {
	if b {
		return 1
	} else {
		return 0
	}
}
