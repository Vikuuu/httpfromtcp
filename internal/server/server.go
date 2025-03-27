package server

import (
	"fmt"
	"log"
	"net"
	"sync/atomic"

	"http/internal/request"
	"http/internal/response"
)

type Server struct {
	handler  Handler
	closed   atomic.Bool
	listener net.Listener
}

func Serve(port int, handler Handler) (*Server, error) {
	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		return nil, err
	}

	s := &Server{
		listener: listener,
		handler:  handler,
	}
	go s.listen()
	return s, nil
}

func (s *Server) Close() error {
	s.closed.Store(true)
	if s.listener != nil {
		return s.listener.Close()
	}
	return nil
}

func (s *Server) listen() {
	for {
		conn, err := s.listener.Accept()
		if err != nil {
			if s.closed.Load() {
				return
			}
			log.Printf("Error accepting conn: %s\n", err)
			continue
		}
		fmt.Println("Connection accepted from: ", conn.RemoteAddr())
		go s.handle(conn)
	}
}

func (s *Server) handle(conn net.Conn) {
	defer fmt.Println("Connection closed with: ", conn.RemoteAddr())
	defer conn.Close()
	r, _ := request.RequestFromReader(conn)
	w := response.NewWriter(conn)
	s.handler(w, r)
}
