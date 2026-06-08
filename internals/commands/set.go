package commands

import (
	"errors"
	"net"
	"redis-demo/internals/db"
	"strconv"
	"strings"
	"time"
)

type SetCommand struct {
	DB *db.DB
}

func (s *SetCommand) Execute(conn net.Conn, args []string) error {
	if len(args) < 2 {
		return errors.New("-ERR wrong number of arguments for 'set' command\r\n")
	}
	var ttl time.Duration = 0

	if len(args) >= 4 && strings.ToUpper(args[2]) == "EX" {

		seconds, err := strconv.Atoi(args[3])

		if err != nil {
			return errors.New("ERR value is not an integer or out of range")
		}

		ttl = time.Duration(seconds) * time.Second
	}
	// Save to database
	s.DB.Set(args[0], args[1], ttl)

	_, err := conn.Write([]byte("+OK\r\n"))
	return err

}
