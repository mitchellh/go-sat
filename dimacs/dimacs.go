// Package dimacs parses the DIMACS CNF format.
//
// DIMACS CNF is a common format used to represent boolean expressions in
// conjunctive normal form. It is often used as a way to test SAT solvers.
//
// This package will only parse the CNF problem type in the file. If the file
// contains any other problem type then parsing will fail even if it is a
// valid syntax otherwise.
//
// The full DIMACS CNF format is explained here:
// http://www.domagoj-babic.com/uploads/ResearchProjects/Spear/dimacs-cnf.pdf
package dimacs

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"strconv"

	"github.com/mitchellh/go-sat/cnf"
)

// Problem is a single problem from a DIMACS CNF file.
//
// Variables and Clauses will reflect exactly what is read from the
// problem line in the DIMACS file. This isn't validated with the number
// of variables or clauses read in the file itself.
type Problem struct {
	Variables int         // Variables is number of declared variables
	Clauses   int         // Clauses is number of declared clauses
	Formula   cnf.Formula // Formula is the actual boolean formula
}

// Parse parses the given input buffer in DIMAC CNF form and returns
// the parsed problem.
func Parse(r io.Reader) (*Problem, error) {
	// Initialize the result
	var result Problem
	result.Variables = -1

	// Create a bufio scanner so we can break it up by line
	scanner := bufio.NewScanner(r)
	read := 0
	for scanner.Scan() {
		raw := scanner.Bytes()

		// Disallow blank lines
		if len(raw) == 0 {
			return nil, fmt.Errorf("blank line not allowed in DIMACS CNF input")
		}

		// If we don't know the number of variables we still haven't
		// seen the problem line. Look for comments or the problem.
		if result.Variables == -1 {
			switch raw[0] {
			case 'c':
				// Ignore, comment line

			case 'p':
				// Problem line found!
				fields := bytes.Fields(raw)
				if len(fields) != 4 {
					return nil, fmt.Errorf(
						"problem line should have 4 fields whitespace separated: %q", raw)
				}

				if string(fields[1]) != "cnf" {
					return nil, fmt.Errorf(
						"problem type must be 'cnf', got: %q", fields[1])
				}

				vars, err := strconv.Atoi(string(fields[2]))
				if err != nil {
					return nil, fmt.Errorf(
						"error converting variable count %q: %s", fields[2], err)
				}

				clauses, err := strconv.Atoi(string(fields[3]))
				if err != nil {
					return nil, fmt.Errorf(
						"error converting clauses count %q: %s", fields[3], err)
				}

				result.Variables = vars
				result.Clauses = clauses

			default:
				return nil, fmt.Errorf(
					"invalid start of line character: %q", raw[0])
			}

			continue
		}

		// Read the line
		fields := bytes.Fields(raw)
		if len(fields) == 0 {
			return nil, fmt.Errorf("clause line should not be empty")
		}

		// The file SHOULD terminate in a zero. However, some CNF files
		// found don't. We accept both.
		if string(fields[len(fields)-1]) == "0" {
			fields = fields[:len(fields)-1]
		}

		// Read all the literals
		c := make([]cnf.Lit, len(fields))
		for i := 0; i < len(fields); i++ {
			val, err := strconv.Atoi(string(fields[i]))
			if err != nil {
				return nil, fmt.Errorf(
					"invalid literal %q", fields[i])
			}

			c[i] = cnf.NewLitInt(val)
		}

		// Add it to our clauses
		result.Formula = append(result.Formula, cnf.Clause(c))

		// Increment our count. If we've read all our expected clauses,
		// then we're done.
		read++
		if read >= result.Clauses {
			break
		}
	}

	return &result, nil
}
