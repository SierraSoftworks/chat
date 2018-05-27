package transports

import (
	"time"

	"github.com/SierraSoftworks/chat/protocol"
)

type Transport interface {
	Active() bool
	Send(cmd *protocol.RawCommand) error
	Receive(timeout time.Duration) (*protocol.RawCommand, error)
	Disconnect() error
}
