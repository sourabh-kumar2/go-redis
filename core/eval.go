package core

import (
	"errors"
	"io"
)

func EvalAndRespond(cmd *RedisCmd, c io.Writer) error {
	switch cmd.Cmd {
	case "PING":
		return evalPING(cmd.Args, c)
	case "SET":
		return evalSET(cmd.Args, c)
	case "GET":
		return evalGET(cmd.Args, c)
	case "TTL":
		return evalTTL(cmd.Args, c)
	case "DEL":
		return evalDel(cmd.Args, c)
	case "EXPIRE":
		return evalExpire(cmd.Args, c)
	default:
		return errors.New("ERR unknown command")
	}
}
