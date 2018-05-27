package server

import (
	"net"

	"github.com/SierraSoftworks/chat/protocol"
	"github.com/SierraSoftworks/chat/transports"
	"github.com/Sirupsen/logrus"
	"github.com/pkg/errors"
)

func newTcpServer(listen string) (*tcpServer, error) {
	l, err := net.Listen("tcp", listen)
	if err != nil {
		return nil, errors.Wrap(err, "tcpServer: failed to start listen server")
	}

	svr := &tcpServer{
		room:     &Room{},
		listener: l,
	}

	return svr, nil
}

type tcpServer struct {
	room     *Room
	listener net.Listener
	shutdown bool
}

func (s *tcpServer) Run() error {
	logrus.Debug("Starting server accept loop")
	for !s.shutdown {
		conn, err := s.listener.Accept()
		if err != nil {
			return err
		}

		logrus.Debug("New client connected")
		tr := transports.NewTcpTransport(conn)
		cl := newClient(s, tr)
		if err := s.room.Join(cl); err != nil {
			return err
		}
	}

	return nil
}

func (s *tcpServer) GetRoom(id string) *Room {
	return s.room
}

func (s *tcpServer) Broadcast(cmd *protocol.RawCommand) error {
	return s.room.Broadcast(cmd)
}

func (s *tcpServer) Shutdown() error {
	return s.listener.Close()
}
