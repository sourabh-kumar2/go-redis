package server

import (
	"syscall"
	"testing"

	"github.com/sourabh-kumar2/go-redis/store"
)

// mockPoller satisfies the Poller interface without touching the OS event queue.
type mockPoller struct {
	removed []int
}

func (p *mockPoller) Add(fd int) error       { return nil }
func (p *mockPoller) Remove(fd int) error    { p.removed = append(p.removed, fd); return nil }
func (p *mockPoller) Wait() ([]Event, error) { return nil, nil }
func (p *mockPoller) Close() error           { return nil }

// socketPair returns (clientFD, serverFD). serverFD is non-blocking, matching
// how the async server configures accepted connections.
func socketPair(t *testing.T) (clientFD, serverFD int) {
	t.Helper()
	fds, err := syscall.Socketpair(syscall.AF_UNIX, syscall.SOCK_STREAM, 0)
	if err != nil {
		t.Fatalf("socketpair: %v", err)
	}
	if err := syscall.SetNonblock(fds[1], true); err != nil {
		syscall.Close(fds[0])
		syscall.Close(fds[1])
		t.Fatalf("setnonblock: %v", err)
	}
	t.Cleanup(func() {
		syscall.Close(fds[0])
		syscall.Close(fds[1])
	})
	return fds[0], fds[1]
}

func TestReadAll(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name  string
		input string
	}{
		{"ping command", "*1\r\n$4\r\nPING\r\n"},
		{"set command", "*3\r\n$3\r\nSET\r\n$3\r\nkey\r\n$5\r\nvalue\r\n"},
		{"no data returns empty", ""},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			clientFD, serverFD := socketPair(t)

			if len(tc.input) > 0 {
				if _, err := syscall.Write(clientFD, []byte(tc.input)); err != nil {
					t.Fatalf("write: %v", err)
				}
			}

			got, closed, err := readAll(serverFD)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if closed {
				t.Fatal("unexpected closed=true")
			}
			if string(got) != tc.input {
				t.Fatalf("got %q, want %q", got, tc.input)
			}
		})
	}
}

func TestReadAllClosedPeer(t *testing.T) {
	t.Parallel()

	clientFD, serverFD := socketPair(t)

	// Close the client end — server should detect the peer closure.
	syscall.Close(clientFD)

	_, closed, err := readAll(serverFD)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !closed {
		t.Fatal("expected closed=true when peer is gone")
	}
}

func TestHandleClient(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name     string
		input    string
		wantResp string
	}{
		{
			name:     "PING returns PONG",
			input:    "*1\r\n$4\r\nPING\r\n",
			wantResp: "+PONG\r\n",
		},
		{
			name:     "PING with message echoes bulk string",
			input:    "*2\r\n$4\r\nPING\r\n$5\r\nhello\r\n",
			wantResp: "$5\r\nhello\r\n",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			clientFD, serverFD := socketPair(t)

			if _, err := syscall.Write(clientFD, []byte(tc.input)); err != nil {
				t.Fatalf("write: %v", err)
			}

			clients := map[int]*fdConn{serverFD: {fd: serverFD}}
			handleClient(serverFD, &mockPoller{}, clients, store.New())

			buf := make([]byte, 512)
			n, err := syscall.Read(clientFD, buf)
			if err != nil {
				t.Fatalf("read response: %v", err)
			}
			if got := string(buf[:n]); got != tc.wantResp {
				t.Fatalf("got %q, want %q", got, tc.wantResp)
			}
		})
	}
}

func TestHandleClientUnknownFD(t *testing.T) {
	t.Parallel()

	_, serverFD := socketPair(t)

	// fd not in clients map — handleClient must return silently.
	handleClient(serverFD, &mockPoller{}, map[int]*fdConn{}, store.New())
}
