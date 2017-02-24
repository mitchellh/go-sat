package sat

// lit represents a literal in a formula.
//
// The least significant bit is the state of the literal
// (0 = positive, 1 = negative). The actual value of the literal is
// the literal right shifted by 1.
//
// Example: for a literal 12, value is  11000  (24)
// Example: for a literal -12, value is 11001 (25)
type lit int

// newLit creates a new literal for the variable v. v must be 0..N and
// s should be true if the variable is negative.
func newLit(v int, s bool) lit { return lit(v + v + boolToInt(s)) }

// Sign reads the sign of the literal. This returns true if the literal is negative.
func (l lit) Sign() bool { return l&1 == 1 }

// Var returns the 0..N variable value for this lit.
func (l lit) Var() int { return int(l >> 1) }

func boolToInt(b bool) int {
	if b {
		return 1
	} else {
		return 0
	}
}
