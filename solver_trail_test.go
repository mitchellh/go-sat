package sat

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/mitchellh/go-sat/cnf"
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
			s.assertLiteral(cnf.NewLitInt(tc.Assert), nil)

			l := cnf.NewLitInt(tc.Lit)
			result := s.ValueLit(l)
			if result != tc.Result {
				t.Fatalf("bad: %s", result)
			}
		})
	}
}

func TestSolverTrimToDecisionLevel(t *testing.T) {
	cases := []struct {
		Assert []int // negative will be decision
		Level  int
		Result []int
	}{
		{
			[]int{-1, -2, 3},
			2,
			[]int{-1, -2, 3},
		},

		{
			[]int{-1, -2, 3},
			1,
			[]int{-1},
		},

		{
			[]int{-1, -2, 3, -4, 5},
			2,
			[]int{-1, -2, 3},
		},
	}

	for i, tc := range cases {
		t.Run(fmt.Sprintf("%d", i), func(t *testing.T) {
			s := New()
			for _, l := range tc.Assert {
				if l < 0 {
					s.newDecisionLevel()
				}

				s.assertLiteral(cnf.NewLitInt(l), nil)
			}

			s.trimToDecisionLevel(tc.Level)

			var result []int
			for _, l := range s.trail {
				result = append(result, l.Int())
			}

			if !reflect.DeepEqual(result, tc.Result) {
				t.Fatalf("bad: %#v", result)
			}
			if s.decisionLevel() != tc.Level {
				t.Fatalf("bad: %d", s.decisionLevel())
			}
		})
	}
}
