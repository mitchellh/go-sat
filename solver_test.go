package sat

import (
	"fmt"
	"testing"
)

func TestSolve(t *testing.T) {
	cases := []struct {
		Name    string
		Formula [][]int
		Result  bool
	}{
		{
			"empty",
			[][]int{},
			true,
		},

		{
			"single literal",
			[][]int{[]int{4}},
			true,
		},
	}

	for i, tc := range cases {
		t.Run(fmt.Sprintf("%d-%s", i, tc.Name), func(t *testing.T) {
			actual := Solve(NewFormulaFromInts(tc.Formula))
			if actual != tc.Result {
				t.Fatalf("bad: %#v", actual)
			}
		})
	}
}
