package atomic

// Atomic is an interface for manging atomic values.
type Atomic interface {
	// Increment increments an atomic value.
	Increment()
	// Decrement decrements an atomic value.
	Decrement()
	// Get gets an atomic value in type of uint32.
	Get() uint32
}
