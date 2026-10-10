package response

import (
	"fmt"

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

func GetDefaultHeaders(contentLen int) headers.Headers {
	h := headers.NewHeaders()
	h.Set("Content-Length", fmt.Sprintf("%d", contentLen))
	h.Set("Connection", "close")
	h.Set("Content-Type", "text/plain")

	return h
}
