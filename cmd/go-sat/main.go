package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/mitchellh/go-sat"
	"github.com/mitchellh/go-sat/dimacs"
)

func main() {
	os.Exit(realMain())
}

func realMain() int {
	flag.Usage = flagUsage
	flag.Parse()

	// Verify args
	args := flag.Args()
	if len(args) != 1 {
		flagUsage()
		return 1
	}

	// Parse the CNF file
	f, err := os.Open(args[0])
	if err != nil {
		printError(err)
		return 1
	}

	p, err := dimacs.Parse(f)
	f.Close()
	if err != nil {
		printError(fmt.Errorf("error parsing cnf file: %s", err))
		return 1
	}

	// Solve the problem
	var s sat.Solver
	s.AddFormula(p.Formula)
	result := s.Solve()
	fmt.Printf("SAT: %v\n", result)
	return 0
}

func flagUsage() {
	fmt.Fprintf(os.Stderr, "Usage: %[1]s [options] <cnf-file>\n", os.Args[0])
	flag.PrintDefaults()
}

func printError(err error) {
	fmt.Fprintf(os.Stderr, err.Error()+"\n")
}
