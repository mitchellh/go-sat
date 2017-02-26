package dimacs

import (
	"fmt"
	"reflect"
	"strings"
	"testing"
)

func TestParse(t *testing.T) {
	cases := []struct {
		Name   string
		Input  string
		Err    bool
		Result [][]int
	}{
		{
			"empty line",
			`
c Foo
p cnf 1 1
1 0
`,
			true,
			nil,
		},

		{
			"basic",
			`c Foo
p cnf 1 1
1 0
`,
			false,
			[][]int{
				[]int{1},
			},
		},

		{
			"multi-clause",
			`c Foo
p cnf 4 3
1 -3 0
2 1 3 0
-4 -2 0
`,
			false,
			[][]int{
				[]int{1, -3},
				[]int{2, 1, 3},
				[]int{-4, -2},
			},
		},

		{
			"end in garbage",
			`c Foo
p cnf 1 1
1 0
%
blank
whatever
`,
			false,
			[][]int{
				[]int{1},
			},
		},
	}

	for i, tc := range cases {
		t.Run(fmt.Sprintf("%d-%s", i, tc.Name), func(t *testing.T) {
			result, err := Parse(strings.NewReader(tc.Input))
			if (err != nil) != tc.Err {
				t.Fatalf("bad: %s", err)
			}
			if err != nil {
				return
			}

			actual := result.Formula.Int()
			if !reflect.DeepEqual(actual, tc.Result) {
				t.Fatalf("bad: %#v", actual)
			}
		})
	}
}
