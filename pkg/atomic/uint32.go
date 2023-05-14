package atomic

import "sync/atomic"

// Value ...
type Value uint32

// NewValue creates and returns new Value.
func NewValue() Value {
	return Value(0)
}

// Increment increments an atomic value.
func (v *Value) Increment() {
	atomic.AddUint32((*uint32)(v), 1)
}

// Decrement decrements an atomic value.
func (v *Value) Decrement() {
	atomic.AddUint32((*uint32)(v), ^uint32(0))
}

// Get gets an atomic value in type of uint32.
func (v *Value) Get() uint32 {
	return uint32(*v)
}
