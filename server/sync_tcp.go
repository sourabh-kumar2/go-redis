package server

import (
	"io"
	"log"
	"net"
	"strconv"

	"github.com/sourabh-kumar2/go-redis/config"
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

			log.Println("command", cmd)
			if err = respond(cmd, c); err != nil {
				log.Println("err write:", err)
			}
		}
	}

}

func readCommand(c net.Conn) (string, error) {
	// TODO: max read at one time is 512 bytes
	// to allow input > 512 bytes, then repeated read until
	// we get EOF or designated delimiter

	buf := make([]byte, 512)
	n, err := c.Read(buf)
	if err != nil {
		return "", err
	}
	return string(buf[:n]), nil
}

func respond(cmd string, c net.Conn) error {
	if _, err := c.Write([]byte(cmd)); err != nil {
		return err
	}
	return nil
}
