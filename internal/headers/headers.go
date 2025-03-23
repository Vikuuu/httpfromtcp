package headers

import (
	"bytes"
	"errors"
	"log"
	"strings"
)

const crlf = "\r\n"

type Headers map[string]string

func NewHeaders() Headers {
	return map[string]string{}
}

func (h Headers) Parse(data []byte) (n int, done bool, err error) {
	log.Printf("[DEBUG] The input data in h.Parse: %q\n", data)
	idx := bytes.Index(data, []byte(crlf))
	if idx == -1 {
		log.Println("[DEBUG] No CRLF found, waiting for full data")
		return 0, false, nil
	}
	if idx == 0 {
		log.Println("[DEBUG] End of headers detected")
		return 2, true, nil
	}

	parts := bytes.SplitN(data[:idx], []byte(":"), 2)
	log.Printf("[DEBUG] The result of spliting data: %q\n", parts)

	if len(parts) < 2 {
		log.Println("[DEBUG] Malformed header")
		return 0, false, errors.New("invalid header")
	}

	key := string(parts[0])
	if key != strings.TrimRight(key, " ") {
		return 0, false, errors.New("invalid header name")
	}

	value := bytes.TrimSpace(parts[1])
	key = strings.TrimSpace(key)
	if !checkKey(key) {
		log.Printf("[DEBUG] Invalid header name: %s\n", key)
		return 0, false, errors.New("invalid header name character")
	}
	log.Printf("[DEBUG] The key-value pairs: %s-%s\n", key, value)

	h.Set(key, string(value))
	log.Printf("[DEBUG] The Return values: %d %t", idx+2, false)
	return idx + 2, false, nil
}

func (h Headers) Set(key, value string) {
	key = strings.ToLower(key)
	if v, ok := h[key]; ok {
		value = v + ", " + value
	}
	h[key] = value
}
