package resp

import (
	"io"
	"log"
)

type Writer struct {
	writer io.Writer
}

func NewWriter(w io.Writer) *Writer {
	return &Writer{writer: w}
}

func (w *Writer) Write(v Value) error {
	var bytes = v.Marshal()
	log.Printf("Responding with: %#v \n", string(bytes[:]))

	_, err := w.writer.Write(bytes)
	if err != nil {
		return err
	}

	return nil
}
