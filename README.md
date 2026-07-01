# go-redis

A Redis-compatible in-memory key-value store written in Go, built from scratch for learning purposes. It implements the [RESP (Redis Serialization Protocol)](https://redis.io/docs/latest/develop/reference/protocol-spec/) and uses a non-blocking async TCP server backed by `kqueue` (macOS) and `epoll` (Linux).

## Features

- **Async TCP server** — event-driven I/O via `kqueue`/`epoll`, supports up to 20,000 concurrent clients
- **RESP protocol** — compatible with `redis-cli` and any standard Redis client
- **Key expiry** — per-key TTL support via `SET ... EX` and `EXPIRE`, with passive (on-read) and active (background cron) eviction of expired keys
- **Key limit & eviction** — store is capped at 1,048,576 keys (`config.KeysLimit`); evicts a key on every `SET` that would exceed the limit
- **Mutex-protected store** — `sync.RWMutex` guards all store reads and writes for safe concurrent access

## Supported Commands

| Command | Syntax | Description |
|---------|--------|-------------|
| `PING` | `PING [message]` | Returns `PONG` or echoes the message |
| `SET` | `SET key value [EX seconds]` | Sets a key with an optional TTL |
| `GET` | `GET key` | Gets the value of a key |
| `DEL` | `DEL key [key ...]` | Deletes one or more keys |
| `TTL` | `TTL key` | Returns remaining TTL in seconds (-1 if no expiry, -2 if missing) |
| `EXPIRE` | `EXPIRE key seconds` | Sets a TTL on an existing key |

## Project Structure

```
.
├── main.go              # Entry point; parses --host and --port flags
├── config/
│   └── config.go        # Host/port/key-limit configuration (default: 0.0.0.0:7379, limit: 1M keys)
├── core/
│   ├── store.go               # In-memory key-value store (RWMutex-protected)
│   ├── eviction.go            # Eviction logic (evicts first key found when store is full)
│   ├── delete_expired_keys.go # Active (background) expired-key eviction via sampling
│   ├── resp.go                # RESP encoder/decoder
│   ├── cmd.go                 # RedisCmd type
│   ├── eval.go                # Command dispatcher
│   └── eval_*.go              # Per-command handlers and tests
└── server/
    ├── async_tcp.go     # Non-blocking TCP server (main server)
    ├── sync_tcp.go      # Blocking TCP server (reference implementation)
    ├── poller.go        # Poller interface
    ├── poller_darwin.go # kqueue implementation (macOS)
    ├── poller_linux.go  # epoll implementation (Linux)
    └── fdconn.go        # File-descriptor-backed net.Conn
```

## Getting Started

**Prerequisites:** Go 1.21+

```bash
# Clone the repository
git clone https://github.com/sourabh-kumar2/go-redis.git
cd go-redis

# Run the server (default: 0.0.0.0:7379)
go run main.go

# Optional flags
go run main.go --host 127.0.0.1 --port 6379
```

Connect with `redis-cli`:

```bash
redis-cli -p 7379

127.0.0.1:7379> PING
PONG
127.0.0.1:7379> SET name Alice EX 60
OK
127.0.0.1:7379> GET name
"Alice"
127.0.0.1:7379> TTL name
(integer) 59
127.0.0.1:7379> DEL name
(integer) 1
```

## Running Tests

```bash
go test ./...
```

## License

[MIT](LICENSE)
