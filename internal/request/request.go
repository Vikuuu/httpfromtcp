package request

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"strings"

	"http/internal/headers"
)

var (
	ErrInvalidMethod      = errors.New("invalid http method")
	ErrInvalidVersion     = errors.New("invalid http version")
	ErrInvalidRequestLine = errors.New("invalid http request line")
	ErrInvalidPath        = errors.New("invalid http path")
)

type state int

const (
	requestStateInitialized state = iota
	requestStateParsingHeaders
	requestStateDone
)

const (
	crlf       = "\r\n"
	bufferSize = 8
)

type Request struct {
	RequestLine RequestLine
	Headers     headers.Headers
	State       state
}

type RequestLine struct {
	HttpVersion   string
	RequestTarget string
	Method        string
}

func RequestFromReader(reader io.Reader) (*Request, error) {
	buf := make([]byte, bufferSize, bufferSize)
	readToIndex := 0
	r := &Request{State: requestStateInitialized, Headers: headers.NewHeaders()}
	for r.State != requestStateDone {
		if readToIndex >= len(buf) {
			tempBuf := make([]byte, 2*len(buf), 2*len(buf))
			copy(tempBuf, buf)
			buf = tempBuf
		}
		nBytesRead, err := reader.Read(buf[readToIndex:])
		if err != nil {
			if errors.Is(err, io.EOF) {
				if r.State != requestStateDone {
					return nil, fmt.Errorf(
						"incomplete request, in state: %d, read n bytes on EOF: %d",
						r.State, nBytesRead)
				}
				break
			}
			return nil, err
		}
		readToIndex += nBytesRead
		nBytesParsed, err := r.parse(buf[:readToIndex])
		if err != nil {
			return nil, err
		}
		copy(buf, buf[nBytesParsed:])
		readToIndex -= nBytesParsed
	}

	return r, nil
}

func parseRequestLine(data []byte) (*RequestLine, int, error) {
	idx := bytes.Index(data, []byte(crlf))
	if idx == -1 {
		// No CRLF found, we need to read more data
		return nil, 0, nil
	}
	requestLineText := string(data[:idx])
	requestLine, err := requestLineFromString(requestLineText)
	if err != nil {
		return nil, 0, err
	}
	return requestLine, idx + 2, nil
}

func requestLineFromString(str string) (*RequestLine, error) {
	parts := strings.Split(str, " ")
	if len(parts) != 3 {
		return nil, ErrInvalidRequestLine
	}
	method, path, ver := parts[0], parts[1], parts[2]
	if !isUpper(method) {
		return nil, ErrInvalidMethod
	}
	if ver != "HTTP/1.1" {
		return nil, ErrInvalidVersion
	}
	ver = strings.Split(ver, "/")[1]
	if !checkPath(path) {
		return nil, ErrInvalidPath
	}

	return &RequestLine{
		Method: method, RequestTarget: path, HttpVersion: ver,
	}, nil
}

func (r *Request) parse(data []byte) (int, error) {
	totalBytesParsed := 0
	for r.State != requestStateDone {
		n, err := r.parseSingle(data[totalBytesParsed:])
		if err != nil {
			return 0, err
		}
		totalBytesParsed += n
		if n == 0 {
			break
		}
	}
	return totalBytesParsed, nil
}

func (r *Request) parseSingle(data []byte) (int, error) {
	switch r.State {
	case requestStateInitialized:
		requestLine, n, err := parseRequestLine(data)
		if err != nil {
			return 0, err
		}
		if n == 0 {
			return 0, nil
		}
		r.RequestLine = *requestLine
		r.State = requestStateParsingHeaders
		return n, nil
	case requestStateParsingHeaders:
		n, done, err := r.Headers.Parse(data)
		if err != nil {
			return 0, err
		}
		if done {
			r.State = requestStateDone
		}
		return n, nil
	case requestStateDone:
		return 0, errors.New("error: trying to read data in a done state")
	default:
		return 0, errors.New("error: unknown state")
	}
}
