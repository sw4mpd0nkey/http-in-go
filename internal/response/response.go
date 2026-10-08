package response

import (
	"fmt"
	"io"

	"boot.swampdonkey.dev/internal/headers"
)

type Response struct {
}

type StatusCode int

const (
	SuccessResponse     StatusCode = 200
	ClientErrorResponse StatusCode = 400
	ServerErrorResponse StatusCode = 500
)

func WriteStatusLine(w io.Writer, statusCode StatusCode) error {

	statusLine := []byte{}

	switch statusCode {

	case SuccessResponse:
		statusLine = []byte("HTTP/1.1 200 OK\r\n")
	case ClientErrorResponse:
		statusLine = []byte("HTTP/1.1 400 Bad Request\r\n")
	case ServerErrorResponse:
		statusLine = []byte("HTTP/1.1 500 Internal Server Error\r\n")
	default:
		return fmt.Errorf("unrecognized error code")
	}

	_, err := w.Write(statusLine)

	return err
}

func GetDefaultHeaders(contentLen int) headers.Headers {
	h := headers.NewHeaders()
	h.Set("Content-Length", fmt.Sprintf("%d", contentLen))
	h.Set("Connection", "close")
	h.Set("Content-Type", "text/plain")

	return h
}

func WriteHeaders(w io.Writer, headers headers.Headers) error {
	b := []byte{}

	headers.ForEach(func(n, v string) {
		b = fmt.Appendf(b, "%s: %s\r\n", n, v)
	})
	b = fmt.Append(b, "\r\n")

	_, err := w.Write(b)
	return err
}
