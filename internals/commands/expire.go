package commands

import (
	"errors"
	"net"
	"redis-demo/internals/db"
	"strconv"
	"time"
)

type ExpireCommand struct {
	DB *db.DB
}

func (e *ExpireCommand) Execute(conn net.Conn, args []string) error {
	if len(args) < 2 {
		return errors.New("ERR wrong number of arguments for 'expire' command")
	}

	seconds, err := strconv.Atoi(args[1])
	if err != nil {
		return errors.New("ERR value is not an integer or out of range")
	}

	result := 0
	if e.DB.Expire(args[0], time.Duration(seconds)*time.Second) {
		result = 1
	}

	resp := []byte(":0\r\n")
	if result == 1 {
		resp = []byte(":1\r\n")
	}

	_, err = conn.Write(resp)
	return err
}
