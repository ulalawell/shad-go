//go:build !solution

package externalsort

import (
	"container/heap"
	"errors"
	"io"
	"os"
)

// mergeItem хранит строку и источник.
type mergeItem struct {
	value  string
	source int
}

// priorityQueue реализует heap.Interface для сортировки строк.
type priorityQueue []mergeItem

func (pq priorityQueue) Len() int            { return len(pq) }
func (pq priorityQueue) Less(i, j int) bool  { return pq[i].value < pq[j].value }
func (pq priorityQueue) Swap(i, j int)       { pq[i], pq[j] = pq[j], pq[i] }
func (pq *priorityQueue) Push(x interface{}) { *pq = append(*pq, x.(mergeItem)) }
func (pq *priorityQueue) Pop() interface{} {
	n := len(*pq)
	item := (*pq)[n-1]
	*pq = (*pq)[:n-1]
	return item
}

// Merge сливает отсортированные входные данные.
func Merge(w LineWriter, readers ...LineReader) error {
	pq := &priorityQueue{}
	heap.Init(pq)

	for i, r := range readers {
		l, err := r.ReadLine()
		if err != nil {
			if errors.Is(err, io.EOF) {
				if l != "" {
					heap.Push(pq, mergeItem{value: l, source: -1})
				}
				continue
			}
			return err
		}
		heap.Push(pq, mergeItem{value: l, source: i})
	}

	for pq.Len() > 0 {
		item := heap.Pop(pq).(mergeItem)
		if err := w.Write(item.value); err != nil {
			return err
		}
		if item.source == -1 {
			continue
		}
		l, err := readers[item.source].ReadLine()
		if err != nil {
			if errors.Is(err, io.EOF) {
				if l != "" {
					heap.Push(pq, mergeItem{value: l, source: -1})
				}
				continue
			}
			return err
		}
		heap.Push(pq, mergeItem{value: l, source: item.source})
	}
	return nil
}

// Sort сортирует содержимое файлов и записывает результат.
func Sort(w io.Writer, in ...string) error {
	var readers []LineReader

	for _, file := range in {
		if err := sortFile(file); err != nil {
			return err
		}

		f, err := os.Open(file)
		if err != nil {
			return err
		}
		defer f.Close()

		readers = append(readers, NewReader(f))
	}

	return Merge(NewWriter(w), readers...)
}

func sortFile(file string) error {
	pq := &priorityQueue{}
	heap.Init(pq)

	if err := readAndPushLines(file, pq); err != nil {
		return err
	}

	return writeSortedLines(file, pq)
}

func readAndPushLines(file string, pq *priorityQueue) error {
	f, err := os.Open(file)
	if err != nil {
		return err
	}
	defer f.Close()

	r := NewReader(f)
	for {
		l, err := r.ReadLine()
		if err != nil {
			if errors.Is(err, io.EOF) {
				if l != "" {
					heap.Push(pq, mergeItem{value: l, source: 0})
				}
				break
			}

			return err
		}
		heap.Push(pq, mergeItem{value: l, source: 0})
	}

	return nil
}

func writeSortedLines(file string, pq *priorityQueue) error {
	f, err := os.Create(file)
	defer f.Close()
	if err != nil {
		return err
	}

	writer := NewWriter(f)
	for pq.Len() > 0 {
		item := heap.Pop(pq).(mergeItem)
		if err = writer.Write(item.value); err != nil {
			return err
		}
	}

	return nil
}
