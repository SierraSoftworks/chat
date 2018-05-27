package transports

import (
	"bufio"
	"context"
	"net"
	"time"

	"github.com/SierraSoftworks/chat/protocol"
	"github.com/pkg/errors"
)

func NewTcpTransport(conn net.Conn) Transport {
	t := &tcpTransport{
		conn:   conn,
		reader: bufio.NewScanner(conn),
		parser: &protocol.Parser{},
	}

	return t
}

type tcpTransport struct {
	conn   net.Conn
	reader *bufio.Scanner
	parser *protocol.Parser
	closed bool
}

func (t *tcpTransport) Active() bool {
	return !t.closed
}

func (t *tcpTransport) Send(cmd *protocol.RawCommand) error {
	_, err := t.conn.Write([]byte(cmd.String() + "\n"))
	if err != nil {
		return errors.Wrap(err, "transport: failed to send")
	}

	return nil
}

func (t *tcpTransport) Receive(timeout time.Duration) (*protocol.RawCommand, error) {
	if err := t.conn.SetReadDeadline(time.Now().Add(timeout)); err != nil {
		return nil, errors.Wrap(err, "transport: failed to set read deadline")
	}

	if !t.reader.Scan() {
		if t.reader.Err() == context.DeadlineExceeded {
			return nil, nil
		}

		return nil, t.reader.Err()
	}

	return t.parser.ParseLine(t.reader.Text())
}

func (t *tcpTransport) Disconnect() error {
	t.closed = true
	return t.conn.Close()
}
