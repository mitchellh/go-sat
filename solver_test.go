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
				Tracer:  &testTracer{T: t},
			}

			actual := s.Solve()
			if actual != tc.Result {
				t.Fatalf("bad: %#v", actual)
			}
		})
	}
}

type testTracer struct {
	T *testing.T
}

func (t *testTracer) Printf(format string, v ...interface{}) {
	t.T.Logf(format, v...)
}
