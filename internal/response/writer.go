package response

import (
	"errors"
	"fmt"
	"io"
	"strconv"

	"http/internal/headers"
)

var ErrOutOfOrder = errors.New("out of order write")

type WriterState int

const (
	writeStateStatusLine WriterState = iota
	writeStateHeader
	writeStateBody
	writeStateDone
)

type Writer struct {
	writer      io.Writer
	writerState WriterState
}

func NewWriter(w io.Writer) *Writer {
	return &Writer{
		writer:      w,
		writerState: writeStateStatusLine,
	}
}

func (w *Writer) WriteHeaders(headers headers.Headers) error {
	if w.writerState != writeStateHeader {
		return ErrOutOfOrder
	}
	for k, v := range headers {
		fieldLine := fmt.Sprintf("%s: %s\r\n", k, v)
		_, err := w.writer.Write([]byte(fieldLine))
		if err != nil {
			return err
		}
	}
	w.writer.Write([]byte("\r\n"))
	w.writerState = writeStateBody
	return nil
}

func (w *Writer) WriteBody(p []byte) (int, error) {
	if w.writerState != writeStateBody {
		return 0, ErrOutOfOrder
	}
	n, err := w.writer.Write(p)
	if err != nil {
		return 0, err
	}
	return n, nil
}

func (w *Writer) WriteChunkedBody(p []byte) (int, error) {
	totalBytes := 0
	dataLenLine := fmt.Sprintf("%s\r\n", strconv.FormatInt(int64(len(p)), 16))
	n, err := w.WriteBody([]byte(dataLenLine))
	if err != nil {
		return 0, err
	}
	totalBytes += n
	n, err = w.WriteBody(p)
	if err != nil {
		return totalBytes, err
	}
	totalBytes += n
	n, err = w.WriteBody([]byte("\r\n"))
	if err != nil {
		return totalBytes, err
	}
	return totalBytes + n, nil
}

func (w *Writer) WriteChunkedBodyDone(h headers.Headers) (int, error) {
	tN := 0
	n, err := w.WriteBody([]byte("0\r\n"))
	if err != nil {
		return 0, err
	}
	tN += n
	w.writerState = writeStateHeader
	err = w.WriteHeaders(h)
	if err != nil {
		return 0, err
	}
	w.writerState = writeStateBody
	n, err = w.WriteBody([]byte("\r\n"))
	if err != nil {
		return tN, nil
	}
	tN += n
	return tN, nil
}
