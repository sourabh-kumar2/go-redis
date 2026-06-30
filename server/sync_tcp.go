package server

import (
	"fmt"
	"io"
	"log"
	"net"
	"strconv"
	"strings"

	"github.com/sourabh-kumar2/go-redis/config"
	"github.com/sourabh-kumar2/go-redis/core"
)

func RunTCPSyncServer() {
	log.Println("starting the synchronous server on", config.Host, config.Port)

	var con_clients uint

	// listening to configured host and port
	lsnr, err := net.Listen("tcp", config.Host+":"+strconv.Itoa(config.Port))
	if err != nil {
		panic(err)
	}
	for {
		// blocking call: waiting for the new client to connect
		c, err := lsnr.Accept()
		if err != nil {
			panic(err)
		}

		con_clients++
		log.Println("client connected with address:", c.RemoteAddr(), "concurrent clients", con_clients)

		for {
			// over the socket: continuously read the command and print it out
			cmd, err := readCommand(c)
			if err != nil {
				c.Close()
				con_clients--

				log.Println("client disconnected:", c.RemoteAddr(), "concurrent clients", con_clients)
				if err == io.EOF {
					break
				}
				log.Println("err", err)
			}

			respond(cmd, c)
		}
	}

}

func readCommand(c net.Conn) (*core.RedisCmd, error) {
	data, err := readAllSync(c)
	if err != nil {
		return nil, err
	}

	tokens, err := core.DecodeArrayString(data)
	if err != nil {
		return nil, err
	}

	return &core.RedisCmd{
		Cmd:  strings.ToUpper(tokens[0]),
		Args: tokens[1:],
	}, nil
}

// readAllSync reads in 512-byte chunks until a short read (the client has
// stopped sending for now) or EOF, so a command isn't truncated at 512 bytes.
func readAllSync(c net.Conn) ([]byte, error) {
	var data []byte
	chunk := make([]byte, 512)
	for {
		n, err := c.Read(chunk)
		if n > 0 {
			data = append(data, chunk[:n]...)
		}
		if err != nil {
			if err == io.EOF && len(data) > 0 {
				return data, nil
			}
			return data, err
		}
		if n < len(chunk) {
			return data, nil
		}
	}
}

func respond(cmd *core.RedisCmd, c net.Conn) {
	err := core.EvalAndRespond(cmd, c)
	if err != nil {
		respondError(err, c)
	}

}

func respondError(err error, c net.Conn) {
	c.Write([]byte(fmt.Sprintf("-%s\r\n", err)))
}
