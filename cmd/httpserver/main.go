package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
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
