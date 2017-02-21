package sat

import (
	"fmt"
	"testing"

	"github.com/mitchellh/go-sat/cnf"
)

func TestTrailIsUnit(t *testing.T) {
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
			m := testTrail(t, tc.Input)
			c := make([]cnf.Literal, len(tc.Clause))
			for i, v := range tc.Clause {
				c[i] = cnf.Literal(v)
			}

			actual := m.IsUnit(cnf.Clause(c), cnf.Literal(tc.Lit))
			if actual != tc.IsUnit {
				t.Fatalf("bad: %v", actual)
			}
		})
	}
}

func testTrail(t *testing.T, in []int) trail {
	var m trail
	for _, v := range in {
		m.Assert(cnf.Literal(v), true)
	}

	return m
}
