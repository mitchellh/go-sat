package sat

import (
	"fmt"
	"reflect"
	"testing"
)

func TestFormulaNegate(t *testing.T) {
	cases := []struct {
		Input  [][]int
		Output [][]int
	}{
		{
			[][]int{
				[]int{1},
				[]int{-3, 4},
			},
			[][]int{
				[]int{-1},
				[]int{3, -4},
			},
		},
	}

	for i, tc := range cases {
		t.Run(fmt.Sprintf("%d", i), func(t *testing.T) {
			actual := NewFormulaFromInts(tc.Input).Negate()
			expected := NewFormulaFromInts(tc.Output)

			if !reflect.DeepEqual(actual, expected) {
				t.Fatalf("bad: %#v", actual)
			}
		})
	}
}
