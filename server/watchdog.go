package server

import (
	"time"

	"github.com/SierraSoftworks/chat/protocol"
)

type watchdog struct {
	client  *Client
	timeout time.Duration
	ticker  *time.Ticker
}

func newWatchdog(cl *Client, interval, timeout time.Duration) *watchdog {
	w := &watchdog{
		client:  cl,
		timeout: timeout,
		ticker:  time.NewTicker(interval),
	}

	go w.pinger()

	return w
}

func (w *watchdog) Stop() {
	w.ticker.Stop()
}

func (w *watchdog) pinger() {
	for range w.ticker.C {
		if err := w.client.Send(protocol.NewPing()); err != nil {
			w.client.Disconnect()
			return
		}

		if time.Now().Sub(w.client.lastSeen) > w.timeout {
			w.client.Disconnect()
			return
		}
	}
}
