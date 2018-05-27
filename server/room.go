package server

import (
	"sync"

	"github.com/Sirupsen/logrus"
	"github.com/pkg/errors"

	"github.com/SierraSoftworks/chat/protocol"
)

type Room struct {
	Clients []*Client

	lock sync.Mutex
}

func (r *Room) Join(cl *Client) error {
	r.lock.Lock()
	defer r.lock.Unlock()

	for _, c := range r.Clients {
		if c == cl {
			return errors.New("client already in the room")
		}
	}

	r.Clients = append(r.Clients, cl)
	logrus.Debug("Client joined the room")
	return nil
}

func (r *Room) Leave(cl *Client) error {
	r.lock.Lock()
	defer r.lock.Unlock()

	for i, c := range r.Clients {
		if c == cl {
			r.Clients = append(r.Clients[:i], r.Clients[i+1:]...)
			logrus.Debug("Client left the room")
			return nil
		}
	}

	return errors.New("client not found in room")
}

func (r *Room) Broadcast(cmd *protocol.RawCommand) error {
	r.lock.Lock()
	defer r.lock.Unlock()
	logrus.WithField("cmd", cmd).Debug("Broadcasting command")

	sendErrors := map[*Client]error{}

	for _, cl := range r.Clients {
		if err := cl.Send(cmd); err != nil {
			sendErrors[cl] = err
		}
	}

	if len(sendErrors) > 0 {
		return errors.New("failed to broadcast to all clients")
	}

	return nil
}
