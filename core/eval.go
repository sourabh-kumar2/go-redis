package core

import (
	"errors"
	"io"
	"log"
	"strconv"
	"time"
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
	default:
		return evalPING(cmd.Args, c)
	}
}

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

func evalSET(args []string, c io.Writer) error {
	if len(args) <= 1 {
		return errors.New("ERR wrong number of arguments for 'set' command")
	}

	exDurationMs := int64(-1)

	key, value := args[0], args[1]

	for i := 2; i < len(args); i++ {
		switch args[i] {
		case "EX", "ex", "Ex", "eX":
			i++
			if i == len(args) {
				return errors.New("ERR syntax error")
			}

			exDurationSec, err := strconv.ParseInt(args[3], 10, 64)
			if err != nil {
				return errors.New("ERR value is not an integer or out of range")
			}

			exDurationMs = exDurationSec * 1000
		default:
			return errors.New("ERR syntax error")
		}
	}

	Put(key, NewObj(value, exDurationMs))

	_, err := c.Write(Encode("OK", true))
	return err
}

func evalGET(args []string, c io.Writer) error {
	if len(args) != 1 {
		return errors.New("ERR wrong number of arguments for 'get' command")
	}

	key := args[0]
	obj := Get(key)
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
