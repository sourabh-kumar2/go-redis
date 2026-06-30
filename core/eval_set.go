package core

import (
	"errors"
	"io"
	"strconv"
)

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
