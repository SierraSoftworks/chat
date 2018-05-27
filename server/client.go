package server

import (
	"fmt"
	"net"
	"time"

	"github.com/SierraSoftworks/chat/protocol"
	"github.com/SierraSoftworks/chat/transports"
	"github.com/Sirupsen/logrus"
	"github.com/pkg/errors"
)

type Client struct {
	transport transports.Transport
	lastSeen  time.Time
	server    Server
	watchdog  *watchdog
}

func newClient(server Server, tr transports.Transport) *Client {
	cl := &Client{
		server:    server,
		transport: tr,
		lastSeen:  time.Now(),
	}

	go cl.pumpMessages()

	cl.watchdog = newWatchdog(cl, 1*time.Second, 5*time.Second)

	return cl
}

func (c *Client) Send(cmd *protocol.RawCommand) error {
	logrus.WithField("cmd", cmd).Debug("Sending command to client")
	return c.transport.Send(cmd)
}

func (c *Client) pumpMessages() {
	for c.transport.Active() {
		cmd, err := c.transport.Receive(5 * time.Second)
		if err != nil {
			if err, ok := err.(net.Error); ok && err.Timeout() {
				continue
			}

			fmt.Println(errors.Wrap(err, "client: failed to receive"))
			continue
		}

		if cmd != nil {
			go c.onCommand(cmd)
		}
	}
}

func (c *Client) Disconnect() error {
	c.transport.Disconnect()
	c.watchdog.Stop()
	return c.server.GetRoom("").Leave(c)
}

func (c *Client) onCommand(cmd *protocol.RawCommand) {
	logrus.WithField("cmd", cmd).Debug("Received command from client")
	switch cmd.Operand {
	case "ping":
		c.lastSeen = time.Now()
	case "msg":
		c.server.GetRoom("").Broadcast(cmd)
	default:
		logrus.WithError(errors.Errorf("client: unknown command type")).Error()
	}
}
