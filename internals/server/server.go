package server

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"strings"

	"redis-demo/internals/commands"
	"redis-demo/internals/protocols"
)

type Server struct {
	addr     string
	listener net.Listener
	registry map[string]commands.Command
}

func NewServer(addr string) *Server {
	s := &Server{
		addr:     addr,
		registry: make(map[string]commands.Command),
	}
	
	s.registry["PING"] = &commands.PingCommand{}

	return s
}

func (s *Server) Start() error {
	var err error
	s.listener, err = net.Listen("tcp", s.addr)
	if err != nil {
		return err
	}
	defer s.listener.Close()

	fmt.Printf("Redis clone server listening on %s\n", s.addr)

	for {
		conn, err := s.listener.Accept()
		if err != nil {
			log.Println("Connection accept error:", err)
			continue
		}
		go s.handleConnection(conn)
	}
}

func (s *Server) handleConnection(conn net.Conn) {
	defer func() {
		fmt.Printf("client disconnected: %v\n", conn.RemoteAddr())
		conn.Close()
	}()

	reader := bufio.NewReader(conn)

	for {
		cmd, err := s.readCommand(reader)
		if err != nil {
			if err == io.EOF {
				return 
			}
			s.respondError(err, conn)
			return
		}

		s.respond(cmd, conn)
	}
}

func (s *Server) readCommand(reader *bufio.Reader) (*commands.RedisCmd, error) {
	
	line, err := reader.ReadBytes('\n')
	if err != nil {
		return nil, err
	}

	if len(line) == 0 || line[0] != '*' {
		return nil, errors.New("ERR protocol error: expected array start '*'")
	}

	
	fullData := append([]byte{}, line...)
	

	count := 0
	for i := 1; i < len(line) && line[i] != '\r'; i++ {
		count = count*10 + int(line[i]-'0')
	}

	
	for i := 0; i < count; i++ {
		
		lenLine, err := reader.ReadBytes('\n')
		if err != nil {
			return nil, err
		}
		fullData = append(fullData, lenLine...)

		
		strLen := 0
		for j := 1; j < len(lenLine) && lenLine[j] != '\r'; j++ {
			strLen = strLen*10 + int(lenLine[j]-'0')
		}

		payload := make([]byte, strLen+2)
		_, err = io.ReadFull(reader, payload)
		if err != nil {
			return nil, err
		}
		fullData = append(fullData, payload...)
	}

	
	tokens, err := protocols.DecodeArrayString(fullData)
	if err != nil {
		return nil, err
	}

	if len(tokens) == 0 {
		return nil, errors.New("empty command")
	}

	return &commands.RedisCmd{
		Cmd:  strings.ToUpper(tokens[0]),
		Args: tokens[1:],
	}, nil
}

func (s *Server) respond(cmd *commands.RedisCmd, conn net.Conn) {
	
	executor, ok := s.registry[cmd.Cmd]
	if !ok {
		s.respondError(errors.New("ERR unknown command"), conn)
		return
	}

	if err := executor.Execute(conn, cmd.Args); err != nil {
		s.respondError(err, conn)
	}
}

func (s *Server) respondError(err error, conn net.Conn) {
	_, _ = conn.Write([]byte(fmt.Sprintf("-%s\r\n", err.Error())))
}