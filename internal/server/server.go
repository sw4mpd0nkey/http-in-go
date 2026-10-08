package server

import (
	"fmt"
	"io"
	"log"
	"net"

	"boot.swampdonkey.dev/internal/response"
)

type Server struct {
	closed bool
}

func runConnection(s *Server, conn io.ReadWriteCloser) {
	defer conn.Close()
	headers := response.GetDefaultHeaders(0)
	response.WriteStatusLine(conn, response.SuccessResponse)
	response.WriteHeaders(conn, headers)

}

func runServer(s *Server, listener net.Listener) error {

	for {
		conn, err := listener.Accept()
		if err != nil {
			return err
		}

		go runConnection(s, conn)
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
