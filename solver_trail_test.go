package sat

import (
	"fmt"
	"testing"

	"github.com/mitchellh/go-sat/cnf"
	"github.com/mitchellh/go-sat/packed"
)

func TestSolverValueLit(t *testing.T) {
	cases := []struct {
		Assert int
		Lit    int
		Result Tribool
	}{
		{
			4,
			-4,
			False,
		},

		{
			4,
			4,
			True,
		},

		{
			-4,
			4,
			False,
		},

		{
			-4,
			-4,
			True,
		},
	}

	for i, tc := range cases {
		t.Run(fmt.Sprintf("%d", i), func(t *testing.T) {
			s := New()
			s.assertLiteral(cnf.Literal(tc.Assert), false)

			l := packed.NewLitInt(tc.Lit)
			result := s.ValueLit(l)
			if result != tc.Result {
				t.Fatalf("bad: %s", result)
			}
		})
	}
}
