//go:build !solution

package otp

import (
	"io"
)

type StreamCipherReader struct {
	inputStream io.Reader
	prng        io.Reader
}

type StreamCipherWriter struct {
	outputStream io.Writer
	prng         io.Reader
}

func NewReader(r io.Reader, prng io.Reader) io.Reader {
	return &StreamCipherReader{
		inputStream: r,
		prng:        prng,
	}
}

func NewWriter(w io.Writer, prng io.Reader) io.Writer {
	return &StreamCipherWriter{
		outputStream: w,
		prng:         prng,
	}
}

func (s *StreamCipherReader) Read(p []byte) (int, error) {
	buf := make([]byte, len(p))

	n, err := s.inputStream.Read(p)
	if n > 0 {
		_, _ = s.prng.Read(buf[:n]) // prng никогда не возвращает ошибку
		for i := 0; i < n; i++ {
			p[i] ^= buf[i]
		}
	}

	return n, err
}

func (s *StreamCipherWriter) Write(p []byte) (int, error) {
	buf := make([]byte, len(p))
	_, _ = s.prng.Read(buf) // prng никогда не возвращает ошибку

	for i := range p {
		buf[i] ^= p[i]
	}

	return s.outputStream.Write(buf)
}
