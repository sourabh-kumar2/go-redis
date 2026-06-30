package core

import (
	"errors"
	"io"
	"log"
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
	log.Println("ttl", obj)

	if obj.ExpiresAt == -1 {
		_, err := c.Write(Encode(-1, false))
		return err
	}

	durationMs := obj.ExpiresAt - time.Now().UnixMilli()

	log.Println("duration", durationMs)

	if durationMs < 0 {
		_, err := c.Write(Encode(-2, false))
		return err
	}
	_, err := c.Write(Encode(durationMs/1000, false))
	return err
}
