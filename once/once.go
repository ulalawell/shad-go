//go:build !solution

package once

type Once struct {
	ch   chan struct{}
	done chan struct{}
}

func New() *Once {
	return &Once{
		ch:   make(chan struct{}, 1),
		done: make(chan struct{}),
	}
}

func (o *Once) Do(f func()) {
	select {
	case o.ch <- struct{}{}:
		defer close(o.done)
		f()
	default:
		<-o.done
	}
}
