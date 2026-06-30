package core

import (
	"errors"
	"io"
)

func evalPING(args []string, c io.Writer) error {
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
