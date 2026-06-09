package commands

import (
	"errors"
	"fmt"
	"net"
	"redis-demo/internals/db"
)

type DelCommand struct {
	DB *db.DB
}

func (d *DelCommand) Execute(conn net.Conn, args []string) error {
	if len(args) < 1 {
		return errors.New("ERR wrong number of arguments for 'del' command")
	}

	deleted := d.DB.Del(args...)

	resp := fmt.Sprintf(":%d\r\n", deleted)
	_, err := conn.Write([]byte(resp))
	return err
}
