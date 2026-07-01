package core

import (
	"errors"
	"io"

	"github.com/sourabh-kumar2/go-redis/store"
)

func evalDel(args []string, c io.Writer, s *store.Store) error {
	if len(args) == 0 {
		return errors.New("ERR wrong number of arguments for 'del' command")
	}
	var countDeleted int
	for _, key := range args {
		if ok := s.Del(key); ok {
			countDeleted++
		}
	}
	_, err := c.Write(Encode(countDeleted, false))
	return err
}
