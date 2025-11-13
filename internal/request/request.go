package request

import (
	"fmt"
	"io"
	"log"
	"strings"
)

type Request struct {
	RequestLine RequestLine
}

type RequestLine struct {
	HttpVersion   string
	RequestTarget string
	Method        string
}

var ERROR_MALFORMED_REQUEST_LINE = fmt.Errorf("malformed request-line")
var ERR_UNSUPPORTED_HTTP_VERSION = fmt.Errorf("http version not supported")
var INCOMPLETE_START_LINE = fmt.Errorf("incomplete start line")
var SEPERATOR = "\r\n"

func parseRequestLine(line string) (*RequestLine, string, error) {
	index := strings.Index(line, SEPERATOR)
	if index == -1 {
		return nil, line, nil
	}

	startLine := line[:index]
	restOfMsg := line[:index+len(SEPERATOR)]

	parts := strings.Split(startLine, " ")
	if len(parts) != 3 {
		return nil, restOfMsg, ERROR_MALFORMED_REQUEST_LINE
	}

	httpParts := strings.Split(parts[2], "/")
	if len(httpParts) != 2 || httpParts[0] != "HTTP" || httpParts[1] != "1.1" {
		return nil, restOfMsg, ERROR_MALFORMED_REQUEST_LINE
	}
	requestline := &RequestLine{
		Method:        parts[0],
		RequestTarget: parts[1],
		HttpVersion:   httpParts[1],
	}

	return requestline, restOfMsg, nil

}
func RequestFromReader(reader io.Reader) (*Request, error) {

	data, err := io.ReadAll(reader)

	if err != nil {
		log.Fatal("error", "error", err)
		panic(err)
	}
	str := string(data)

	rl, _, err := parseRequestLine(str)

	if err != nil {
		return nil, err
	}
	return &Request{
		RequestLine: *rl,
	}, err
}
