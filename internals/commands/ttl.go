package commands

import (
	"errors"
	"fmt"
	"net"
	"redis-demo/internals/db"
)

type TtlCommand struct {
	DB *db.DB
}


func (t *TtlCommand) Execute(conn net.Conn, args []string) error {

	if len(args) < 1 {
		return errors.New("ERR wrong number of arguments for 'ttl' command")
	}

	ttlSec := t.DB.TTL(args[0])

	resp := fmt.Sprintf(":%d\r\n", ttlSec)
	_, err := conn.Write([]byte(resp))
	return err
}
