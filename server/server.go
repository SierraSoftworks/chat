package server

import "github.com/SierraSoftworks/chat/protocol"

type Server interface {
	GetRoom(id string) *Room
	Broadcast(cmd *protocol.RawCommand) error
	Run() error
	Shutdown() error
}
