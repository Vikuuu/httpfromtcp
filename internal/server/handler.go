package server

import (
	"http/internal/request"
	"http/internal/response"
)

type Handler func(w *response.Writer, req *request.Request)

type HandlerError struct {
	StatusCode int
	Message    string
}

func (he *HandlerError) Write(w *response.Writer) error {
	err := w.WriteStatusLine(response.StatusCode(he.StatusCode))
	h := response.GetDefaultHeaders(len(he.Message))
	w.WriteHeaders(h)
	if err != nil {
		return err
	}
	w.WriteBody([]byte(he.Message))
	return nil
}
