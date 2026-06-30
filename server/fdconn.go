package server

import (
	"io"
	"net"
	"syscall"
	"time"
)

// fdConn adapts a raw, poller-managed file descriptor to the net.Conn
// interface so it can be passed straight into core.EvalAndRespond
// without involving Go's runtime netpoller.
type fdConn struct {
	fd int
}

func (c *fdConn) Read(b []byte) (int, error) {
	n, err := syscall.Read(c.fd, b)
	if err != nil {
		return 0, err
	}
	if n == 0 {
		return 0, io.EOF
	}
	return n, nil
}

func (c *fdConn) Write(b []byte) (int, error) {
	return syscall.Write(c.fd, b)
}

func (c *fdConn) Close() error {
	return syscall.Close(c.fd)
}

func (c *fdConn) LocalAddr() net.Addr                { return nil }
func (c *fdConn) RemoteAddr() net.Addr               { return nil }
func (c *fdConn) SetDeadline(t time.Time) error      { return nil }
func (c *fdConn) SetReadDeadline(t time.Time) error  { return nil }
func (c *fdConn) SetWriteDeadline(t time.Time) error { return nil }
