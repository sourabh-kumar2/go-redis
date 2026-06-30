package server

const maxEvents = 256

// Event represents a single ready file descriptor, translated from
// whatever OS-specific event the backend (epoll/kqueue) produced.
type Event struct {
	Fd       int
	Readable bool
	Writable bool
}

// Poller is the OS-agnostic interface for a readiness event queue,
// implemented by poller_linux.go (epoll) and poller_darwin.go (kqueue).
type Poller interface {
	// Add registers fd for read (and optionally write) readiness events.
	Add(fd int) error

	// Remove deregisters fd from the queue, e.g. before closing it.
	Remove(fd int) error

	// Wait blocks until one or more registered fds are ready and
	// returns the events that fired.
	Wait() ([]Event, error)

	// Close releases the underlying queue fd.
	Close() error
}
