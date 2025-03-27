package headers

import (
	"bytes"
	"errors"
	"strings"
)

const crlf = "\r\n"

type Headers map[string]string

func NewHeaders() Headers {
	return map[string]string{}
}

func (h Headers) Parse(data []byte) (n int, done bool, err error) {
	idx := bytes.Index(data, []byte(crlf))
	if idx == -1 {
		return 0, false, nil
	}
	if idx == 0 {
		return 2, true, nil
	}

	parts := bytes.SplitN(data[:idx], []byte(":"), 2)

	if len(parts) < 2 {
		return 0, false, errors.New("invalid header")
	}

	key := string(parts[0])
	if key != strings.TrimRight(key, " ") {
		return 0, false, errors.New("invalid header name")
	}

	value := bytes.TrimSpace(parts[1])
	key = strings.TrimSpace(key)
	if !checkKey(key) {
		return 0, false, errors.New("invalid header name character")
	}

	h.Set(key, string(value))
	return idx + 2, false, nil
}

func (h Headers) Set(key, value string) {
	key = strings.ToLower(key)
	if v, ok := h[key]; ok {
		value = v + ", " + value
	}
	h[key] = value
}

func (h Headers) Override(key, value string) {
	key = strings.ToLower(key)
	h[key] = value
}

func (h Headers) Get(key string) string {
	return h[strings.ToLower(key)]
}

func (h Headers) Remove(key string) {
	delete(h, strings.ToLower(key))
}
