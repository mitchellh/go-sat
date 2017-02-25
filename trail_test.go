package sat

import (
	"fmt"
	"testing"

	"github.com/mitchellh/go-sat/cnf"
)

func TestTrailTrimToLevel(t *testing.T) {
	cases := []struct {
		Name  string
		Input []int
		Level int
		End   []int
	}{
		{
			"simple trim",
			[]int{106, 1, 2, 103, 4},
			1,
			[]int{106, 1, 2},
		},

		{
			"trim to zero",
			[]int{106, 1, 2, 103, 4},
			0,
			[]int{},
		},
	}

	for i, tc := range cases {
		t.Run(fmt.Sprintf("%d-%s", i, tc.Name), func(t *testing.T) {
			m := testTrail(t, tc.Input)
			t.Logf("Input: %s", m)
			m.TrimToLevel(tc.Level)
			expected := testTrail(t, tc.End)
			if m.String() != expected.String() {
				t.Fatalf("bad: %s", m)
			}
		})
	}
}

func testTrail(t *testing.T, in []int) *trail {
	m := newTrail()
	for _, v := range in {
		decision := false
		if v > 100 {
			decision = true
			v -= 100
		}

		m.Assert(cnf.Literal(v), decision)
	}

	return m
}
