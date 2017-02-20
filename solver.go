package sat

// Solve solves the given formula, returning ture on satisfiability and
// false on unsatisfiability. This is just temporary. We'll return the
// actual values for solving eventually.
func Solve(f Formula) bool {
	var m trail

	varsF := f.Vars()
	for {
		// If the trail contains the same number of elements as
		// the variables in the formula, then we've found a satisfaction.
		if len(m) == len(varsF) {
			return true
		}

	}

	return false
}
