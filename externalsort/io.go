//go:build !change

package externalsort

import "io"

type LineReader interface {
	ReadLine() (string, error)
}

type LineWriter interface {
	Write(l string) error
}

type LineReaderTool struct {
	reader   io.Reader
	diffData []byte
}

type LineWriterTool struct {
	writer   io.Writer
	diffData []byte
}

func NewReader(r io.Reader) LineReader {
	return &LineReaderTool{
		reader:   r,
		diffData: []byte{},
	}
}

func NewWriter(w io.Writer) LineWriter {
	return &LineWriterTool{
		writer:   w,
		diffData: []byte{},
	}
}

func (lr *LineReaderTool) ReadLine() (string, error) {
	for {
		if len(lr.diffData) != 0 {
			for i := 0; i < len(lr.diffData); i++ {
				if lr.diffData[i] == '\n' {
					oneLine := string(lr.diffData[:i])
					lr.diffData = lr.diffData[i+1:]
					return oneLine, nil
				}
			}
		}

		p := make([]byte, 10)
		n, err := lr.reader.Read(p)
		p = p[:n]

		for i := 0; i < len(p); i++ {
			if p[i] == '\n' {
				oneLine := string(lr.diffData) + string(p[:i])
				lr.diffData = p[i+1:]
				return oneLine, nil
			}
		}

		if err != nil {
			res := string(lr.diffData) + string(p)
			lr.diffData = []byte{}
			return res, err
		}

		lr.diffData = append(lr.diffData, p...)
	}

}

func (lw *LineWriterTool) Write(l string) error {
	l += "\n"
	p := []byte(l)

	_, err := lw.writer.Write(p)

	return err
}
