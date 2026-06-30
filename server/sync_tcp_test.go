package server

import (
	"fmt"
	"net"
	"testing"

	"github.com/sourabh-kumar2/go-redis/core"
)

// pipe returns a connected client/server net.Conn pair and registers cleanup.
func pipe(t *testing.T) (client, server net.Conn) {
	t.Helper()
	client, server = net.Pipe()
	t.Cleanup(func() { client.Close(); server.Close() })
	return
}

func TestReadAllSync(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name  string
		input string
	}{
		{"ping command", "*1\r\n$4\r\nPING\r\n"},
		{"set command", "*3\r\n$3\r\nSET\r\n$3\r\nfoo\r\n$3\r\nbar\r\n"},
		{"get command", "*2\r\n$3\r\nGET\r\n$3\r\nfoo\r\n"},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			client, server := pipe(t)

			go client.Write([]byte(tc.input))

			got, err := readAllSync(server)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if string(got) != tc.input {
				t.Fatalf("got %q, want %q", got, tc.input)
			}
		})
	}
}

func TestReadCommand(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name     string
		input    string
		wantCmd  string
		wantArgs []string
		wantErr  bool
	}{
		{
			name:     "PING",
			input:    "*1\r\n$4\r\nPING\r\n",
			wantCmd:  "PING",
			wantArgs: []string{},
		},
		{
			name:     "SET with args",
			input:    "*3\r\n$3\r\nSET\r\n$3\r\nfoo\r\n$3\r\nbar\r\n",
			wantCmd:  "SET",
			wantArgs: []string{"foo", "bar"},
		},
		{
			name:     "lowercase cmd is uppercased",
			input:    "*2\r\n$3\r\nget\r\n$3\r\nfoo\r\n",
			wantCmd:  "GET",
			wantArgs: []string{"foo"},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			client, server := pipe(t)

			go client.Write([]byte(tc.input))

			cmd, err := readCommand(server)
			if tc.wantErr {
				if err == nil {
					t.Fatal("expected error, got nil")
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if cmd.Cmd != tc.wantCmd {
				t.Fatalf("cmd: got %q, want %q", cmd.Cmd, tc.wantCmd)
			}
			if fmt.Sprintf("%v", cmd.Args) != fmt.Sprintf("%v", tc.wantArgs) {
				t.Fatalf("args: got %v, want %v", cmd.Args, tc.wantArgs)
			}
		})
	}
}

func TestRespondError(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name   string
		errMsg string
		want   string
	}{
		{"generic error", "ERR something went wrong", "-ERR something went wrong\r\n"},
		{"syntax error", "ERR syntax error", "-ERR syntax error\r\n"},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			client, server := pipe(t)

			done := make(chan string, 1)
			go func() {
				buf := make([]byte, 512)
				n, _ := client.Read(buf)
				done <- string(buf[:n])
			}()

			respondError(fmt.Errorf("%s", tc.errMsg), server)

			if got := <-done; got != tc.want {
				t.Fatalf("got %q, want %q", got, tc.want)
			}
		})
	}
}

func TestRespond(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name     string
		cmd      *core.RedisCmd
		wantResp string
	}{
		{
			name:     "PING returns PONG",
			cmd:      &core.RedisCmd{Cmd: "PING", Args: []string{}},
			wantResp: "+PONG\r\n",
		},
		{
			name:     "PING with message echoes message",
			cmd:      &core.RedisCmd{Cmd: "PING", Args: []string{"hello"}},
			wantResp: "$5\r\nhello\r\n",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			client, server := pipe(t)

			done := make(chan string, 1)
			go func() {
				buf := make([]byte, 512)
				n, _ := client.Read(buf)
				done <- string(buf[:n])
			}()

			respond(tc.cmd, server)

			if got := <-done; got != tc.wantResp {
				t.Fatalf("got %q, want %q", got, tc.wantResp)
			}
		})
	}
}
