package server

import (
	"fmt"
	"net"
)

type Server struct {
	addr string
}

func New(addr string) *Server {
	return &Server{
		addr: addr,
	}
}

func (s *Server) Start() error {
	listener, err := net.Listen("tcp", s.addr)
	if err != nil {
		return err
	}

	defer listener.Close()

	fmt.Printf("Server listening on ", s.addr)

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Printf("accept error: ", err)
			continue
		}

		fmt.Printf("new client: ", conn.RemoteAddr())

		go s.handleConnection(conn)
	}
}

func (s *Server) handleConnection(conn net.Conn) {
	defer func() {
		fmt.Printf("client disconnected: ", conn.RemoteAddr())
		conn.Close()
	}()

	buf := make([]byte, 4096)

	for {
		n, err := conn.Read(buf)
		if err != nil {
			return
		}

		data := buf[:n]


		fmt.Printf("received: ", data)

		_, err = conn.Write([]byte("+OK\r\n"))
		if err != nil {
			return
		}
	}
}