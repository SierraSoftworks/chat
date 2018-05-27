package server

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/pkg/errors"
	"gopkg.in/urfave/cli.v1"
)

var Command = cli.Command{
	Name:  "server",
	Usage: "Run a chat server to which clients may connect",
	Flags: []cli.Flag{
		cli.StringFlag{
			Name:  "listen",
			Usage: "The address on which the server should listen",
			Value: "0.0.0.0:2428",
		},
	},
	Action: func(c *cli.Context) error {
		svr, err := newTcpServer(c.String("listen"))
		if err != nil {
			return errors.Wrap(err, "server: failed to start listening")
		}

		shutdownCh := make(chan os.Signal)
		signal.Notify(shutdownCh, syscall.SIGINT, syscall.SIGTERM)

		go func() {
			select {
			case <-shutdownCh:
				svr.Shutdown()
			}
		}()

		return svr.Run()
	},
}
