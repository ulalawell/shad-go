//go:build !solution

package externalsort

import (
	"bufio"
	"container/heap"
	"io"
	"os"
)

type StringHeap []string

func (h StringHeap) Len() int           { return len(h) }
func (h StringHeap) Less(i, j int) bool { return h[i] < h[j] }
func (h StringHeap) Swap(i, j int)      { h[i], h[j] = h[j], h[i] }

func (h *StringHeap) Push(x any) {
	// Push and Pop use pointer receivers because they modify the slice's length,
	// not just its contents.
	*h = append(*h, x.(string))
}

func (h *StringHeap) Pop() any {
	old := *h
	n := len(old)
	x := old[n-1]
	*h = old[0 : n-1]
	return x
}

func Merge(w LineWriter, readers ...LineReader) error {
	h := &StringHeap{}
	heap.Init(h)

	for i := 0; i < len(readers); i++ {
		for {
			line, err := readers[i].ReadLine()
			if err == nil || line != "" {
				heap.Push(h, line)
			}

			if err == io.EOF {
				break
			} else if err != nil {
				return err
			}
		}
	}

	for len(*h) > 0 {
		err := w.Write(heap.Pop(h).(string))
		if err != nil {
			return err
		}
	}

	return nil
}

func Sort(w io.Writer, in ...string) error {
	readers := []LineReader{}

	for _, file := range in {
		f, err := os.Open(file)

		if err != nil {
			return err
		}

		readers = append(readers, NewReader(bufio.NewReader(f)))
		defer f.Close()
	}

	lw := NewWriter(w)

	err := Merge(lw, readers...)

	return err
}
