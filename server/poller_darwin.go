//go:build darwin

package server

import "syscall"

type kqueuePoller struct {
	fd int
}

func NewPoller() (Poller, error) {
	fd, err := syscall.Kqueue()
	if err != nil {
		return nil, err
	}
	return &kqueuePoller{fd: fd}, nil
}

func (p *kqueuePoller) Add(fd int) error {
	changes := []syscall.Kevent_t{
		{
			Ident:  uint64(fd),
			Filter: syscall.EVFILT_READ,
			Flags:  syscall.EV_ADD,
		},
	}
	_, err := syscall.Kevent(p.fd, changes, nil, nil)
	return err
}

func (p *kqueuePoller) Remove(fd int) error {

	changes := []syscall.Kevent_t{
		{
			Ident:  uint64(fd),
			Filter: syscall.EVFILT_READ,
			Flags:  syscall.EV_DELETE,
		},
	}
	_, err := syscall.Kevent(p.fd, changes, nil, nil)
	return err
}

func (p *kqueuePoller) Wait() ([]Event, error) {
	raw := make([]syscall.Kevent_t, maxEvents)
	n, err := syscall.Kevent(p.fd, nil, raw, nil)
	if err != nil {
		if err == syscall.EINTR {
			return nil, nil
		}
		return nil, err
	}

	events := make([]Event, 0, n)
	for i := 0; i < n; i++ {
		events = append(events, Event{
			Fd:       int(raw[i].Ident),
			Readable: raw[i].Filter == syscall.EVFILT_READ,
			Writable: raw[i].Filter == syscall.EVFILT_WRITE,
		})
	}
	return events, nil
}

func (p *kqueuePoller) Close() error {
	return syscall.Close(p.fd)
}
