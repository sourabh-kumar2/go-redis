package core

import (
	"errors"
	"io"
	"time"
)

func evalTTL(args []string, c io.Writer) error {
	if len(args) != 1 {
		return errors.New("ERR wrong number of arguments for 'get' command")
	}

	key := args[0]
	obj := Get(key)

	if obj == nil {
		_, err := c.Write(Encode(-2, false))
		return err
	}

	if obj.ExpiresAt == -1 {
		_, err := c.Write(Encode(-1, false))
		return err
	}

	durationMs := obj.ExpiresAt - time.Now().UnixMilli()

	if durationMs < 0 {
		_, err := c.Write(Encode(-2, false))
		return err
	}
	_, err := c.Write(Encode(durationMs/1000, false))
	return err
}
