package sat

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/mitchellh/go-sat/dimacs"
)

func BenchmarkSolver_satLib(b *testing.B) {
	// Get the dirs containing our tests, this will be sorted already
	dirs := satlibDirs(b)

	// Go through each dir and run the tests
	for _, d := range dirs {
		// Run the tests for this dir
		b.Run(filepath.Base(d), func(b *testing.B) {
			satlibBenchmarkDir(b, d)
		})
	}
}

func satlibBenchmarkDir(b *testing.B, dir string) {
	// Open the directory so we can read each file
	dirF, err := os.Open(dir)
	if err != nil {
		b.Fatalf("err: %s", err)
	}

	entries, err := dirF.Readdirnames(-1)
	dirF.Close()
	if err != nil {
		b.Fatalf("err: %s", err)
	}

	// Go through each entry and attempt to solve each
	count := 0
	for _, entry := range entries {
		// Ignore non-CNF files
		if filepath.Ext(entry) != ".cnf" {
			continue
		}

		// Test this entry
		b.Run(entry, func(b *testing.B) {
			satlibBenchmarkFile(b, filepath.Join(dir, entry))
		})

		// Run only the threshold number
		count++
		if count >= satlibBenchThreshold {
			return
		}
	}
}

func satlibBenchmarkFile(b *testing.B, path string) {
	// Parse the problem
	f, err := os.Open(path)
	if err != nil {
		b.Fatalf("err: %s", err)
	}

	p, err := dimacs.Parse(f)
	f.Close()
	if err != nil {
		b.Fatalf("err: %s", err)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		s := New()
		s.AddFormula(p.Formula.Pack())
		s.Solve()
	}
}
