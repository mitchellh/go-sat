package cnf

import (
	"fmt"
	"testing"
)

func TestLit(t *testing.T) {
	cases := []struct {
		Input int
		Var   int
		Sign  bool
	}{
		{
			12,
			12,
			false,
		},

		{
			-12,
			12,
			true,
		},
	}

	for i, tc := range cases {
		t.Run(fmt.Sprintf("%d", i), func(t *testing.T) {
			v := tc.Input
			sign := v < 0
			if v < 0 {
				v *= -1
			}

			lit := NewLit(v, sign)
			if lit.Var() != tc.Var {
				t.Fatalf("bad: %d", lit.Var())
			}
			if lit.Sign() != tc.Sign {
				t.Fatalf("bad: %v", lit.Sign())
			}

			// Test negation
			neg := lit.Neg()
			if neg.Var() != tc.Var {
				t.Fatalf("bad: %d", lit.Var())
			}
			if neg.Sign() == tc.Sign {
				t.Fatalf("bad: %v", neg.Sign())
			}
		})
	}
}
