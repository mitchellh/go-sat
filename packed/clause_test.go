package packed

import (
	"fmt"
	"testing"
)

func TestClauseSetLits_maxVar(t *testing.T) {
	cases := []struct {
		Input  []int
		MaxVar int
	}{
		{
			[]int{1, -3, -12, 4},
			12,
		},
	}

	for i, tc := range cases {
		t.Run(fmt.Sprintf("%d", i), func(t *testing.T) {
			c := NewClause(0)
			var ls []Lit
			for _, l := range tc.Input {
				ls = append(ls, NewLitInt(l))
			}
			c.SetLits(ls)

			if v := c.MaxVar(); v != tc.MaxVar {
				t.Fatalf("bad: %d", v)
			}
		})
	}
}
