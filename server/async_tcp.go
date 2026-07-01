package server

import (
	"log"
	"net"
	"strings"
	"syscall"
	"time"

	"github.com/sourabh-kumar2/go-redis/config"
	"github.com/sourabh-kumar2/go-redis/core"
	"github.com/sourabh-kumar2/go-redis/store"
)

var con_clients uint = 0
var cronFrequency time.Duration = 1 * time.Second
var lastCronExecTime time.Time = time.Now()

func RunTCPAsyncServer() error {
	log.Println("starting the asynchronous server on", config.Host, config.Port)

	max_clients := 20_000

	s := store.New()

	// create a socket
	serverFD, err := syscall.Socket(syscall.AF_INET, syscall.SOCK_STREAM, 0)
	if err != nil {
		return err
	}
	defer syscall.Close(serverFD)

	// Set the socket to operate in non-blocking mode.
	if err = syscall.SetNonblock(serverFD, true); err != nil {
		return err
	}

	// Bind IP and port
	ip4 := net.ParseIP(config.Host)
	if err := syscall.Bind(serverFD, &syscall.SockaddrInet4{
		Port: config.Port,
		Addr: [4]byte{ip4[0], ip4[1], ip4[2], ip4[3]},
	}); err != nil {
		return err
	}

	// Start listening
	if err = syscall.Listen(serverFD, max_clients); err != nil {
		return err
	}

	poller, err := NewPoller()
	if err != nil {
		return err
	}
	defer poller.Close()

	if err = poller.Add(serverFD); err != nil {
		return err
	}

	clients := make(map[int]*fdConn)

	for {
		if time.Now().After(lastCronExecTime.Add(cronFrequency)) {
			s.DeleteExpiredKeys()
			lastCronExecTime = time.Now()
		}

		events, err := poller.Wait()
		if err != nil {
			return err
		}

		for _, ev := range events {
			if ev.Fd == serverFD {
				acceptClients(serverFD, poller, clients)
				continue
			}
			if ev.Readable {
				handleClient(ev.Fd, poller, clients, s)
			}
		}
	}
}

// acceptClients drains every pending connection on the listening socket
// in one go, since a single readiness event can represent several
// queued connections.
func acceptClients(serverFD int, poller Poller, clients map[int]*fdConn) {
	for {
		connFD, _, err := syscall.Accept(serverFD)
		if err != nil {
			if err == syscall.EAGAIN || err == syscall.EWOULDBLOCK {
				return
			}
			log.Println("accept error:", err)
			return
		}

		if err := syscall.SetNonblock(connFD, true); err != nil {
			syscall.Close(connFD)
			continue
		}

		if err := poller.Add(connFD); err != nil {
			syscall.Close(connFD)
			continue
		}

		clients[connFD] = &fdConn{fd: connFD}
		con_clients++
		log.Println("client connected, fd", connFD, "concurrent clients", con_clients)
	}
}

func handleClient(fd int, poller Poller, clients map[int]*fdConn, s *store.Store) {
	conn, ok := clients[fd]
	if !ok {
		return
	}

	data, closed, err := readAll(fd)
	if err != nil {
		closeClient(fd, poller, clients)
		return
	}
	if closed {
		closeClient(fd, poller, clients)
		return
	}
	if len(data) == 0 {
		return
	}

	tokens, err := core.DecodeArrayString(data)
	if err != nil {
		log.Println("decode error:", err)
		return
	}
	if len(tokens) == 0 {
		return
	}

	cmd := &core.RedisCmd{
		Cmd:  strings.ToUpper(tokens[0]),
		Args: tokens[1:],
	}

	if err := core.EvalAndRespond(cmd, conn, s); err != nil {
		conn.Write([]byte("-" + err.Error() + "\r\n"))
	}
}

// readAll drains fd in 512-byte chunks until the kernel has no more data
// ready (EAGAIN) or the peer has closed the connection, so callers aren't
// limited to a single 512-byte read per readiness event.
func readAll(fd int) (data []byte, closed bool, err error) {
	chunk := make([]byte, 512)
	for {
		n, err := syscall.Read(fd, chunk)
		if err != nil {
			if err == syscall.EAGAIN || err == syscall.EWOULDBLOCK {
				return data, false, nil
			}
			return nil, false, err
		}
		if n == 0 {
			return data, true, nil
		}

		data = append(data, chunk[:n]...)
		if n < len(chunk) {
			// short read: fd is very likely drained for now.
			return data, false, nil
		}
	}
}

func closeClient(fd int, poller Poller, clients map[int]*fdConn) {
	poller.Remove(fd)
	syscall.Close(fd)
	delete(clients, fd)
	con_clients--
	log.Println("client disconnected, fd", fd, "concurrent clients", con_clients)
}
