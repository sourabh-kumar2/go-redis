package core

import (
	"errors"
	"net"
)

func EvalAndRespond(cmd *RedisCmd, c net.Conn) error {
	switch cmd.Cmd {
	case "PING":
		return evalPING(cmd.Args, c)
	default:
		return evalPING(cmd.Args, c)
	}
}

func evalPING(args []string, c net.Conn) error {
	var b []byte

	switch len(args) {
	case 0:
		b = Encode("PONG", true)
	case 1:
		b = Encode(args[0], false)
	default:
		return errors.New("ERR wrong number of arguments for 'ping' command")
	}

	_, err := c.Write(b)
	return err
}
