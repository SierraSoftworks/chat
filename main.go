package main

import (
	"github.com/SierraSoftworks/chat/client"
	"github.com/SierraSoftworks/chat/server"
	"github.com/Sirupsen/logrus"
	"gopkg.in/urfave/cli.v1"
)

func main() {
	app := cli.NewApp()

	app.Commands = []cli.Command{
		server.Command,
		client.Command,
	}

	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "log-level",
			Usage: "Set the log level used to display information",
			Value: "INFO",
		},
	}

	app.Before = func(c *cli.Context) error {
		switch c.String("log-level") {
		case "ERROR":
			logrus.SetLevel(logrus.ErrorLevel)
		case "WARN":
			logrus.SetLevel(logrus.WarnLevel)
		case "DEBUG":
			logrus.SetLevel(logrus.DebugLevel)
		default:
			logrus.SetLevel(logrus.InfoLevel)
		}

		return nil
	}

	app.RunAndExitOnError()
}
