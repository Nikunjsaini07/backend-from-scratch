package commands

import (
	"fmt"
	"net"
)

type PingCommand struct{}

func (p *PingCommand) Execute(conn net.Conn, args []string) error {
	if len(args) == 0 {
		_, err := conn.Write([]byte("+PONG\r\n"))
		return err
	}

	message := args[0]
	resp := fmt.Sprintf("$%d\r\n%s\r\n", len(message), message)
	_, err := conn.Write([]byte(resp))
	return err
}