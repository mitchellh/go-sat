package sat

import (
	"fmt"
	"reflect"
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
			s.assertLiteral(packed.NewLitInt(tc.Assert))

			l := packed.NewLitInt(tc.Lit)
			result := s.ValueLit(l)
			if result != tc.Result {
				t.Fatalf("bad: %s", result)
			}
		})
	}
}

func TestSolverIsUnit(t *testing.T) {
	cases := []struct {
		Name   string
		Input  []int
		Clause []int
		Lit    int
		IsUnit bool
	}{
		{
			"single element",
			[]int{},
			[]int{4},
			4,
			true,
		},

		{
			"single element trail contains",
			[]int{4},
			[]int{4},
			4,
			false,
		},

		{
			"two element unit",
			[]int{3},
			[]int{-3, 4},
			4,
			true,
		},

		{
			"two element unit with negative",
			[]int{3},
			[]int{-3, 4},
			-4,
			true,
		},

		{
			"two element non-unit",
			[]int{-3},
			[]int{-3, 4},
			4,
			false,
		},

		{
			"three element unit",
			[]int{1, 3, -6},
			[]int{-1, -3, 5},
			5,
			true,
		},
	}

	for i, tc := range cases {
		t.Run(fmt.Sprintf("%d-%s", i, tc.Name), func(t *testing.T) {
			s := New()
			for _, l := range tc.Input {
				s.assertLiteral(packed.NewLitInt(l))
			}

			c := make([]cnf.Literal, len(tc.Clause))
			for i, v := range tc.Clause {
				c[i] = cnf.Literal(v)
			}

			actual := s.isUnit(cnf.Clause(c), cnf.Literal(tc.Lit))
			if actual != tc.IsUnit {
				t.Fatalf("bad: %v", actual)
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

				s.assertLiteral(packed.NewLitInt(l))
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
