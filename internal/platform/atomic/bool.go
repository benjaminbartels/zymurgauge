package atomic

import "sync/atomic"

// Bool is a atomic boolean value
type Bool struct{ flag int32 }

// Set sets the boolean
func (b *Bool) Set(value bool) {
	var i int32
	if value {
		i = 1
	}
	atomic.StoreInt32(&(b.flag), i)
}

// Get returns the boolean
func (b *Bool) Get() bool {
	return atomic.LoadInt32(&(b.flag)) != 0
}
