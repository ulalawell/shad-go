//go:build !change

package externalsort

import (
	"bufio"
	"io"
	"strings"
)

type LineReader interface {
	ReadLine() (string, error)
}

type LineWriter interface {
	Write(l string) error
}

type lineReader struct {
	buffReader *bufio.Reader
}

func NewReader(r io.Reader) LineReader {
	return &lineReader{
		buffReader: bufio.NewReader(r),
	}
}

func (r *lineReader) ReadLine() (string, error) {
	st, err := r.buffReader.ReadString('\n')
	return strings.Trim(st, "\n"), err
}

type lineWriter struct {
	buffWriter *bufio.Writer
}

func NewWriter(w io.Writer) LineWriter {
	return &lineWriter{
		buffWriter: bufio.NewWriter(w),
	}
}

func (w *lineWriter) Write(l string) error {
	_, err := w.buffWriter.WriteString(l + "\n")
	if err != nil {
		return err
	}
	return w.buffWriter.Flush()
}
