package sat

import (
	"fmt"
	"testing"

	"github.com/mitchellh/go-sat/cnf"
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
			s := &Solver{
				Formula: cnf.NewFormulaFromInts(tc.Formula),
				Trace:   true,
			}

			actual := s.Solve()
			if actual != tc.Result {
				t.Fatalf("bad: %#v", actual)
			}
		})
	}
}
