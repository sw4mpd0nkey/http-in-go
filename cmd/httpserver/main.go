package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"boot.swampdonkey.dev/internal/request"
	"boot.swampdonkey.dev/internal/response"
	"boot.swampdonkey.dev/internal/server"
)

const port = 42069

func ReadFile(fileName string) []byte {

	data, err := os.ReadFile("html/" + fileName)
	if err != nil {
		panic(err)
	}

	return data
}

func main() {

	server, err := server.Serve(port, func(w *response.Writer, req *request.Request) *server.HandlerError {

		h := response.GetDefaultHeaders(0)
		body := ReadFile("200.html")
		status := response.SuccessResponse

		if req.RequestLine.RequestTarget == "/yourproblem" {

			body = ReadFile("400.html")
			status = response.ClientErrorResponse

		} else if req.RequestLine.RequestTarget == "/myproblem" {

			body = ReadFile("500.html")
			status = response.ServerErrorResponse

		} else if strings.HasPrefix(req.RequestLine.RequestTarget, "/httpbin/stream") {

			target := req.RequestLine.RequestTarget
			url := "https://httpbingo.org/" + target[len("/httpbin/"):]
			resp, err := http.Get(url)

			if err != nil {
				body = ReadFile("500.html")
				status = response.ServerErrorResponse
			} else {

				w.WriteStatusLine(response.SuccessResponse)

				h.Delete("Content-Length")
				h.Set("Transfer-Encoding", "chunked")
				h.Replace("Content-Type", "text/plain")
				w.WriteHeaders(h)

				for {
					data := make([]byte, 32)
					n, err := resp.Body.Read(data)
					if err != nil {
						break
					}
					w.WriteChunkedBody(data[:n])
				}
				w.WriteChunkedBodyDone()
				return nil
			}
		}

		w.WriteStatusLine(status)
		h.Replace("Content-Length", fmt.Sprintf("%d", len(body)))
		h.Replace("Content-Type", "text/html writer ")
		w.WriteHeaders(h)
		w.WriteBody(body)

		return nil
	})

	if err != nil {
		log.Fatalf("Error starting server: %v", err)
	}

	defer server.Close()
	log.Println("Server started on port", port)

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan
	log.Println("Server gracefully stopped")
}
