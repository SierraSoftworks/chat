package client

import (
	"fmt"
	"net"
	"time"

	"github.com/Sirupsen/logrus"

	"github.com/SierraSoftworks/chat/protocol"
	"github.com/SierraSoftworks/chat/transports"
	"github.com/pkg/errors"
)

func NewClient(name string, tr transports.Transport) *Client {
	cl := &Client{
		Name:      name,
		transport: tr,
		messages:  make(chan protocol.MessageCommand),
	}

	go cl.pumpMessages()

	return cl
}

type Client struct {
	Name      string
	transport transports.Transport
	messages  chan protocol.MessageCommand
}

func (c *Client) SendMessage(msg string) error {
	return c.transport.Send(protocol.NewMessage(c.Name, msg))
}

func (c *Client) Messages() <-chan protocol.MessageCommand {
	return c.messages
}

func (c *Client) Disconnect() error {
	return c.transport.Disconnect()
}

func (c *Client) pumpMessages() {
	for c.transport.Active() {
		logrus.Debug("Receiving new messages from server")
		cmd, err := c.transport.Receive(5 * time.Second)
		if err != nil {
			if err, ok := err.(net.Error); ok && err.Timeout() {
				continue
			}

			fmt.Println(errors.Wrap(err, "client: failed to receive"))
			continue
		}

		if cmd != nil {
			logrus.WithField("cmd", cmd).Debug("Received message from server")
			go c.onCommand(cmd)
		}
	}
}

func (c *Client) onCommand(cmd *protocol.RawCommand) {
	switch cmd.Operand {
	case "ping":
		if err := c.transport.Send(cmd); err != nil {
			if err, ok := err.(net.Error); ok && err.Timeout() {
				logrus.WithField("cmd", cmd).Warn("Timeout during send")
			}
		}
	case "msg":
		c.messages <- protocol.MessageCommand(*cmd)
	default:
		logrus.WithError(errors.Errorf("client: unknown command type")).Error()
	}
}
