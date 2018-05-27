package transports

import (
	"context"
	"time"

	"github.com/SierraSoftworks/chat/protocol"
)

func NewLoopbackTransport() Transport {
	tr := &channelTransport{
		ch: make(chan *protocol.RawCommand, 1),
	}
	tr.peer = tr

	return tr
}

func NewTestTransports() (Transport, Transport) {
	tr1 := &channelTransport{
		ch: make(chan *protocol.RawCommand, 1),
	}

	tr2 := &channelTransport{
		ch: make(chan *protocol.RawCommand, 1),
	}

	tr1.peer, tr2.peer = tr2, tr1

	return tr1, tr2
}

type channelTransport struct {
	peer *channelTransport
	ch   chan *protocol.RawCommand
}

func (t *channelTransport) Active() bool {
	return t.peer != nil
}

func (t *channelTransport) Send(cmd *protocol.RawCommand) error {
	select {
	case t.peer.ch <- cmd:
		return nil
	case <-time.After(100 * time.Millisecond):
		return context.DeadlineExceeded
	}
}

func (t *channelTransport) Receive(timeout time.Duration) (*protocol.RawCommand, error) {
	select {
	case cmd := <-t.ch:
		return cmd, nil
	case <-time.After(timeout):
		return nil, nil
	}
}

func (t *channelTransport) Disconnect() error {
	t.peer = nil
	return nil
}
