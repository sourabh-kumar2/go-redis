package core

import (
	"errors"
	"io"
	"time"

	"github.com/sourabh-kumar2/go-redis/store"
)

func evalGET(args []string, c io.Writer, s *store.Store) error {
	if len(args) != 1 {
		return errors.New("ERR wrong number of arguments for 'get' command")
	}

	key := args[0]
	obj := s.Get(key)
	if obj == nil {
		_, err := c.Write([]byte(RESP_NIL))
		return err
	}

	if obj.ExpiresAt != -1 && obj.ExpiresAt <= time.Now().UnixMilli() {
		_, err := c.Write([]byte(RESP_NIL))
		return err
	}
	_, err := c.Write(Encode(obj.Value, false))
	return err
}
