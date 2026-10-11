package main

import (
	"crypto/sha256"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"boot.swampdonkey.dev/internal/headers"
	"boot.swampdonkey.dev/internal/request"
	"boot.swampdonkey.dev/internal/response"
	"boot.swampdonkey.dev/internal/server"
)

const port = 42069

func ReadFile(fileName string) []byte {

	data, err := os.ReadFile(fileName)
	if err != nil {
		panic(err)
	}

	return data
}

func main() {

	server, err := server.Serve(port, handler)

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

func handler(w *response.Writer, req *request.Request) *server.HandlerError {

	if req.RequestLine.RequestTarget == "/yourproblem" {
		handler400(w, req)
		return nil
	} else if req.RequestLine.RequestTarget == "/myproblem" {
		handler500(w, req)
		return nil
	} else if strings.HasPrefix(req.RequestLine.RequestTarget, "/httpbin/") {
		proxyHandler(w, req)
		return nil
	} else if req.RequestLine.RequestTarget == "/video" {
		videoHandler(w, req)
		return nil
	}

	handler200(w, req)
	return nil
}

func videoHandler(w *response.Writer, _ *request.Request) {

	body := ReadFile("assets/vim.mp4")
	h := response.GetDefaultHeaders(len(body))
	h.Replace("Content-Type", "video/mp4")

	w.WriteStatusLine(response.SuccessResponse)
	w.WriteHeaders(h)
	w.WriteBody(body)
}

func handler200(w *response.Writer, _ *request.Request) {

	body := ReadFile("html/200.html")
	h := response.GetDefaultHeaders(len(body))
	h.Replace("Content-Type", "text/html")

	w.WriteStatusLine(response.SuccessResponse)
	w.WriteHeaders(h)
	w.WriteBody(body)
}

func handler400(w *response.Writer, _ *request.Request) {

	body := ReadFile("html/400.html")
	h := response.GetDefaultHeaders(len(body))
	h.Replace("Content-Type", "text/html")

	w.WriteStatusLine(response.SuccessResponse)
	w.WriteHeaders(h)
	w.WriteBody(body)
}

func handler500(w *response.Writer, _ *request.Request) {

	body := ReadFile("html/500.html")
	h := response.GetDefaultHeaders(len(body))
	h.Replace("Content-Type", "text/html")

	w.WriteStatusLine(response.SuccessResponse)
	w.WriteHeaders(h)
	w.WriteBody(body)
}

func proxyHandler(w *response.Writer, req *request.Request) {

	target := req.RequestLine.RequestTarget
	url := "https://httpbingo.org/" + target[len("/httpbin/"):]
	fmt.Println("Proxying to", url)
	resp, err := http.Get(url)

	if err != nil {
		handler500(w, req)
		return
	}
	defer resp.Body.Close()

	// first write the status line
	w.WriteStatusLine(response.SuccessResponse)

	// next write the headers
	h := response.GetDefaultHeaders(0)
	h.Delete("Content-Length")
	h.Replace("Transfer-Encoding", "chunked")
	h.Replace("Content-Type", "text/plain")
	h.Set("Trailer", "X-Content-Sha256")
	h.Set("Trailer", "X-Content-Length")
	w.WriteHeaders(h)

	// next write the body in chunked segements
	fullBody := make([]byte, 0)
	maxChunkSize := 1024

	for {
		data := make([]byte, maxChunkSize)
		n, err := resp.Body.Read(data)

		if err != nil {
			break
		}

		if n > 0 {
			fullBody = append(fullBody, data[:n]...)
			w.WriteChunkedBody(data[:n])
		}

	}
	w.WriteChunkedBodyDone()

	// writing trailers
	trailer := headers.NewHeaders()
	sha256 := fmt.Sprintf("%x", sha256.Sum256(fullBody))
	length := fmt.Sprintf("%d", len(fullBody))
	trailer.Replace("X-Content-Sha256", sha256)
	trailer.Replace("X-Content-Length", length)
	err = w.WriteTrailers(trailer)

	if err != nil {
		fmt.Println("Error writing trailers:", err)
	}
}
