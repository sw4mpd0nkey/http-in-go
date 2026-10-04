package server

import (
	"fmt"
	"io"
	"log"
	"net"
)

type Server struct {
	closed bool
}

func runConnection(conn io.ReadWriteCloser) {

	out := []byte("HTTP/1.1 200 OK\r\nContent-Type: text/plain\r\nContent-Length: 13\r\n\r\nHello World!")
	conn.Write(out)
	conn.Close()
}

func runServer(s *Server, listener net.Listener) error {

	for {
		conn, err := listener.Accept()
		if err != nil {
			return err
		}

		go runConnection(conn)
	}

}

func Serve(port uint16) (*Server, error) {

	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", port))

	if err != nil {
		log.Fatal("error", "error", err)
		return nil, err
	}

	server := &Server{
		closed: false,
	}

	go runServer(server, listener)

	return server, nil
}

func (s *Server) Close() error {
	s.closed = true
	return nil
}

func (s *Server) Listen() {

}

func (s *Server) Handle(conn net.Conn) {

}
