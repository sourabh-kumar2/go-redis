package core

import (
	"errors"
	"io"
	"strconv"
	"time"
)

func evalExpire(args []string, c io.Writer) error {
	if len(args) <= 1 {
		return errors.New("ERR wrong number of arguments for 'expire' command")
	}
	key := args[0]
	durationSec, err := strconv.ParseInt(args[1], 10, 64)
	if err != nil {
		return errors.New("ERR value is not an integer or out of range")
	}

	obj := Get(key)
	if obj == nil {
		_, err = c.Write(Encode(0, false))
		return err
	}

	obj.ExpiresAt = time.Now().UnixMilli() + durationSec*1000
	_, err = c.Write(Encode(1, false))
	return err
}
