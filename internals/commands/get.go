package commands

import (
	"errors"
	"fmt"
	"net"
	"redis-demo/internals/db"
)

type GetCommand struct {
	DB *db.DB
}

func (g *GetCommand) Execute(conn net.Conn, args []string) error {

	if len(args) < 1 {
		return errors.New("ERR wrong number of arguments for 'get' command")
	}

	val, exists := g.DB.Get(args[0])

	if !exists {
		_, err := conn.Write([]byte("$-1\r\n"))
		return err
	}

	resp := fmt.Sprintf("$%d\r\n%s\r\n", len(val), val)
	_, err := conn.Write([]byte(resp))
	return err
}
