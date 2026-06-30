package core

import (
	"errors"
	"io"
)

func evalDel(args []string, c io.Writer) error {
	if len(args) == 0 {
		return errors.New("ERR wrong number of arguments for 'del' command")
	}
	var countDeleted int
	for _, key := range args {
		if ok := Del(key); ok {
			countDeleted++
		}
	}
	_, err := c.Write(Encode(countDeleted, false))
	return err
}
