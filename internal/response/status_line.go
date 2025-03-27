package response

import "fmt"

type StatusCode int

const (
	_                                    = iota
	StatusOk                  StatusCode = 200
	StatusBadRequest                     = 400
	StatusInternalServerError            = 500
)

func (w *Writer) WriteStatusLine(statusCode StatusCode) error {
	if w.writerState != writeStateStatusLine {
		return ErrOutOfOrder
	}
	phrase := ""
	switch statusCode {
	case StatusOk:
		phrase = "OK"
	case StatusBadRequest:
		phrase = "Bad Request"
	case StatusInternalServerError:
		phrase = "Internal Server Error"
	default:
		phrase = ""
	}
	response := fmt.Sprintf("HTTP/1.1 %d %s\r\n", statusCode, phrase)
	_, err := w.writer.Write([]byte(response))
	if err != nil {
		return err
	}
	w.writerState = writeStateHeader
	return nil
}
