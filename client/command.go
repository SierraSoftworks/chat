package client

import (
	"bufio"
	"fmt"
	"net"
	"os"

	"github.com/SierraSoftworks/chat/transports"
	"github.com/Sirupsen/logrus"
	"github.com/pkg/errors"
	"gopkg.in/urfave/cli.v1"
)

var Command = cli.Command{
	Name:  "client",
	Usage: "Run a chat client which will connect to a server",
	Flags: []cli.Flag{
		cli.StringFlag{
			Name:  "server",
			Usage: "The address of the server to connect to",
			Value: "localhost:2428",
		},
	},
	Before: func(c *cli.Context) error {
		if c.NArg() < 1 {
			return errors.Errorf("you must provide a name")
		}

		return nil
	},
	Action: func(c *cli.Context) error {
		conn, err := net.Dial("tcp", c.String("server"))
		if err != nil {
			return errors.Wrap(err, "chat: failed to connect to server")
		}

		tr := transports.NewTcpTransport(conn)
		cl := NewClient(c.Args().First(), tr)
		defer cl.Disconnect()

		go func() {
			logrus.Debug("Reading messages from the client")
			for msg := range cl.Messages() {
				logrus.Debug("Printing message from the client")
				fmt.Printf("[%s]: %s\n", msg.Sender(), msg.Message())
			}
			logrus.Debug("Stopped reading messages from the client")
		}()

		scanner := bufio.NewScanner(os.Stdin)
		logrus.Debug("")
		for scanner.Scan() {
			cl.SendMessage(scanner.Text())
		}

		return errors.Wrap(scanner.Err(), "chat: failed to read message")
	},
}
