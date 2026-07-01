package core

import (
	"errors"
	"io"

	"github.com/sourabh-kumar2/go-redis/store"
)

func EvalAndRespond(cmd *RedisCmd, c io.Writer, s *store.Store) error {
	switch cmd.Cmd {
	case "PING":
		return evalPING(cmd.Args, c)
	case "SET":
		return evalSET(cmd.Args, c, s)
	case "GET":
		return evalGET(cmd.Args, c, s)
	case "TTL":
		return evalTTL(cmd.Args, c, s)
	case "DEL":
		return evalDel(cmd.Args, c, s)
	case "EXPIRE":
		return evalExpire(cmd.Args, c, s)
	default:
		return errors.New("ERR unknown command")
	}
}
