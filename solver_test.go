package sat

import (
	"fmt"
	"testing"
)

func TestSolve(t *testing.T) {
	Trace = true

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

		{
			"unsatisfiable with backtrack",
			[][]int{
				[]int{4},
				[]int{6},
				[]int{-4, -6},
			},
			false,
		},

		{
			"satisfiable with backtrack",
			[][]int{
				[]int{-4},
				[]int{4, -6},
			},
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
