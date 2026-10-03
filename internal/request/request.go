package request

import (
	"bytes"
	"fmt"
	"io"

	"boot.swampdonkey.dev/internal/headers"
)

const (
	StateInit    parserState = "init"
	StateDone    parserState = "done"
	StateHeaders parserState = "headers"
	StateError   parserState = "error"
)

var ErrMalformedRequestLine = fmt.Errorf("malformed request-line")
var ErrUnsupportedHttpVersion = fmt.Errorf("http version not supported")
var ErrRequeestInErrorState = fmt.Errorf("request in err state")
var ErrIncompleteStartLine = fmt.Errorf("incomplete start line")
var SEPERATOR = []byte("\r\n")

type Request struct {
	RequestLine RequestLine
	Headers     headers.Headers
	state       parserState
}

func (r *Request) parse(data []byte) (int, error) {

	read := 0
outer:
	for {
		currentData := data[read:]

		switch r.state {

		case StateError:
			return 0, ErrRequeestInErrorState

		case StateInit:
			rl, n, err := parseRequestLine(currentData)
			if err != nil {
				r.state = StateError
				return 0, err
			}
			if n == 0 {
				break outer
			}

			r.RequestLine = *rl
			read += n
			r.state = StateHeaders

		case StateHeaders:
			n, done, err := r.Headers.Parse(currentData)

			if err != nil {
				r.state = StateError
				return 0, err
			}

			if n == 0 {
				break outer
			}

			read += n

			if done {
				r.state = StateDone
			}

		case StateDone:
			break outer

		default:
			panic("We done fucked up somehow")
		}

	}
	return read, nil
}

func (r *Request) done() bool {
	return r.state == StateDone
}

func (r *Request) error() bool {
	return r.state == StateError
}

type RequestLine struct {
	HttpVersion   string
	RequestTarget string
	Method        string
}

type parserState string

func newRequest() *Request {
	return &Request{
		state:   StateInit,
		Headers: headers.NewHeaders(),
	}
}

func parseRequestLine(b []byte) (*RequestLine, int, error) {
	index := bytes.Index(b, SEPERATOR)
	if index == -1 {
		return nil, 0, nil
	}

	startLine := b[:index]
	read := index + len(SEPERATOR)

	parts := bytes.Split(startLine, []byte(" "))
	if len(parts) != 3 {
		return nil, 0, ErrMalformedRequestLine
	}

	httpParts := bytes.Split(parts[2], []byte("/"))
	if len(httpParts) != 2 || string(httpParts[0]) != "HTTP" || string(httpParts[1]) != "1.1" {
		return nil, 0, ErrMalformedRequestLine
	}
	requestline := &RequestLine{
		Method:        string(parts[0]),
		RequestTarget: string(parts[1]),
		HttpVersion:   string(httpParts[1]),
	}

	return requestline, read, nil

}

func RequestFromReader(reader io.Reader) (*Request, error) {

	// data, err := io.ReadAll(reader)
	request := newRequest()

	buf := make([]byte, 1024)
	bufIdx := 0

	for !request.done() && !request.error() {
		n, err := reader.Read(buf[bufIdx:])
		if err != nil {
			return nil, err
		}
		bufIdx += n
		readN, err := request.parse(buf[:bufIdx])
		if err != nil {
			return nil, err
		}
		copy(buf, buf[readN:bufIdx])
		bufIdx -= readN

	}
	return request, nil
}
