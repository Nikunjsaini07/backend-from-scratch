package commands

import "net"


type RedisCmd struct {
	Cmd  string
	Args []string
}


type Command interface {
	Execute(conn net.Conn, args []string) error
}