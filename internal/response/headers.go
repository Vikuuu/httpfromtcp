package response

import (
	"fmt"

	"http/internal/headers"
)

func GetDefaultHeaders(contentLen int) headers.Headers {
	h := headers.NewHeaders()
	h.Set("Content-Length", fmt.Sprintf("%d", contentLen))
	h.Set("Connection", "close")
	h.Set("Content-Type", "text/plain")

	return h
}

func (w *Writer) WriteTrailers(h headers.Headers) error {
	err := w.WriteHeaders(h)
	if err != nil {
		return err
	}
	return nil
}
