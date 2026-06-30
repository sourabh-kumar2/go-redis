//go:build linux

package server

import "syscall"

type epollPoller struct {
	fd int
}

func NewPoller() (Poller, error) {
	fd, err := syscall.EpollCreate1(syscall.EPOLL_CLOEXEC)
	if err != nil {
		return nil, err
	}
	return &epollPoller{fd: fd}, nil
}

func (p *epollPoller) Add(fd int) error {
	event := syscall.EpollEvent{
		Events: syscall.EPOLLIN,
		Fd:     int32(fd),
	}
	return syscall.EpollCtl(p.fd, syscall.EPOLL_CTL_ADD, fd, &event)
}

func (p *epollPoller) Remove(fd int) error {
	return syscall.EpollCtl(p.fd, syscall.EPOLL_CTL_DEL, fd, nil)
}

func (p *epollPoller) Wait() ([]Event, error) {
	raw := make([]syscall.EpollEvent, maxEvents)
	n, err := syscall.EpollWait(p.fd, raw, -1)
	if err != nil {
		if err == syscall.EINTR {
			return nil, nil
		}
		return nil, err
	}

	events := make([]Event, 0, n)
	for i := 0; i < n; i++ {
		events = append(events, Event{
			Fd:       int(raw[i].Fd),
			Readable: raw[i].Events&syscall.EPOLLIN != 0,
			Writable: raw[i].Events&syscall.EPOLLOUT != 0,
		})
	}
	return events, nil
}

func (p *epollPoller) Close() error {
	return syscall.Close(p.fd)
}
