package sat

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"testing"

	"github.com/mitchellh/go-sat/cnf"
	"github.com/mitchellh/go-sat/dimacs"
)

var flagSatlib = flag.Bool("satlib", false, "run ALL SATLIB tests (slow!)")

// satlibThreshold is the number of satlib tests to run per category
// when flagSatlib is NOT set. This can be increased as the efficiency of
// the solver improves.
const satlibThreshold = 1

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

// Test the solver with SATLIB problems.
func TestSolver_satlib(t *testing.T) {
	// Get the dirs containing our tests, this will be sorted already
	dirs := satlibDirs(t)

	// Go through each dir and run the tests
	for _, d := range dirs {
		// Run the tests for this dir
		t.Run(filepath.Base(d), func(t *testing.T) {
			satlibTestDir(t, d)
		})
	}
}

func satlibTestDir(t *testing.T, dir string) {
	base := filepath.Base(dir)

	// If the directory has the prefix "sat" then we expect all
	// tests within to be satisfiable. If not, we expect the opposite.
	sat := strings.HasPrefix(base, "sat-")

	// Open the directory so we can read each file
	dirF, err := os.Open(dir)
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	entries, err := dirF.Readdirnames(-1)
	dirF.Close()
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	// Go through each entry and attempt to solve each
	for i, entry := range entries {
		// Ignore non-CNF files
		if filepath.Ext(entry) != ".cnf" {
			continue
		}

		// If we're not running all satlib tests, then step at 3 examples each
		if !*flagSatlib && i >= satlibThreshold {
			break
		}

		// Test this entry
		t.Run(entry, func(t *testing.T) {
			satlibTestFile(t, filepath.Join(dir, entry), sat)
		})
	}
}

func satlibTestFile(t *testing.T, path string, expected bool) {
	// Parse the problem
	f, err := os.Open(path)
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	p, err := dimacs.Parse(f)
	f.Close()
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	// Solve it
	s := &Solver{
		Formula: p.Formula,
	}

	actual := s.Solve()
	if actual != expected {
		t.Fatalf("expected %v, got %v", expected, actual)
	}
}

func satlibDirs(t *testing.T) []string {
	base := filepath.Join("testdata", "satlib")
	dir, err := os.Open(base)
	defer dir.Close()
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	entries, err := dir.Readdir(-1)
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	// Go through and get all the directories (which contain the tests)
	tests := make([]string, 0, len(entries))
	for _, e := range entries {
		if e.IsDir() {
			tests = append(tests, e.Name())
		}
	}

	// Sort so we have a predictable ordering
	sort.Strings(tests)

	// Create the full path
	for i, t := range tests {
		tests[i] = filepath.Join(base, t)
	}

	return tests
}

// testTracer is a Tracer implementation that sends output to the test logger.
type testTracer struct {
	T *testing.T
}

func (t *testTracer) Printf(format string, v ...interface{}) {
	t.T.Logf(format, v...)
}
