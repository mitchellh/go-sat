// Package packed contains a "packed formula" representation. This
// representation is the final representation that the go-sat solver uses
// to do SAT solving.
//
// It is optimized for that SAT solver and isn't meant to be generally
// useful or have the best API experience compared to packages such as
// CNF.
package packed
