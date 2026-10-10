package response

import (
	"fmt"
	"io"

	"boot.swampdonkey.dev/internal/headers"
)

type writerState int

const (
	writerStateStatusLine writerState = iota
	writerStateHeaders
	writerStateBody
	writerStateTrailers
)

type Writer struct {
	writer      io.Writer
	writerState writerState
}

func NewWriter(writer io.Writer) *Writer {
	return &Writer{
		writer:      writer,
		writerState: writerStateStatusLine,
	}
}

func (w *Writer) WriteStatusLine(statusCode StatusCode) error {

	if w.writerState != writerStateStatusLine {
		return fmt.Errorf("cannot write status line in state %d", w.writerState)
	}
	defer func() { w.writerState = writerStateHeaders }()
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

	_, err := w.writer.Write(statusLine)

	return err
}

func (w *Writer) WriteHeaders(headers headers.Headers) error {

	if w.writerState != writerStateHeaders {
		return fmt.Errorf("cannot write headers in state %d", w.writerState)
	}
	defer func() { w.writerState = writerStateBody }()

	b := []byte{}

	headers.ForEach(func(n, v string) {
		b = fmt.Appendf(b, "%s: %s\r\n", n, v)
	})
	b = fmt.Append(b, "\r\n")

	_, err := w.writer.Write(b)
	return err
}

func (w *Writer) WriteBody(p []byte) (int, error) {
	if w.writerState != writerStateBody {
		return 0, fmt.Errorf("cannot write body in state %d", w.writerState)
	}

	n, err := w.writer.Write(p)

	return n, err
}

func (w *Writer) WriteChunkedBody(p []byte) (int, error) {
	if w.writerState != writerStateBody {
		return 0, fmt.Errorf("cannot write body in state %d", w.writerState)
	}

	// start by writing hex
	w.WriteBody([]byte(fmt.Sprintf("%x\r\n", len(p))))
	n, err := w.WriteBody(p)
	w.WriteBody([]byte("\r\n"))

	return n, err
}

func (w *Writer) WriteChunkedBodyDone() (int, error) {
	if w.writerState != writerStateBody {
		return 0, fmt.Errorf("cannot write body in state %d", w.writerState)
	}

	n, err := w.writer.Write([]byte("0\r\n"))
	w.writerState = writerStateTrailers
	return n, err
}

func (w *Writer) WriteTrailers(h headers.Headers) error {
	if w.writerState != writerStateTrailers {
		return fmt.Errorf("cannot write trailers in state %d", w.writerState)
	}
	defer func() { w.writerState = writerStateBody }()
	for k, v := range h {
		_, err := w.writer.Write([]byte(fmt.Sprintf("%s: %s\r\n", k, v)))
		if err != nil {
			return err
		}
	}
	_, err := w.writer.Write([]byte("\r\n"))
	return err
}
