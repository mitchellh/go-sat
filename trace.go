package sat

// Tracer is the implementation for a tracer to use with Solver.
//
// The Go stdlib Logger implements this interface. However, this interface
// lets you slide in anything that adheres to this making it simple to use
// a different log package.
type Tracer interface {
	Printf(format string, v ...interface{})
}
