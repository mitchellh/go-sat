package sat

// Tracer is the implementation for a tracer to use with Solver.
type Tracer interface {
	Printf(format string, v ...interface{})
}
