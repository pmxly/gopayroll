package main

import (
	"gopayroll/common"
	"gopayroll/db"
	boogooServer "gopayroll/server"
	"github.com/urfave/cli"
	"os"
)

var (
	workApp        *cli.App
)

func init() {
	// Initialise a CLI app
	workApp = cli.NewApp()
	workApp.Name = "Boogoo Process Server"
	workApp.Usage = "Process Boogoo Tasks"
	workApp.Author = "Boogoo Team"
	workApp.Email = "tech@boogoo.cn"
	workApp.Version = "1.0.0"
	workApp.Flags = []cli.Flag{
		cli.StringFlag{
			Name:        "c",
			Value:       "",
			Destination: &common.ConfigPath,
			Usage:       "Path to a configuration file",
		},
	}
}

func main() {
	// Set the CLI app commands
	workApp.Commands = []cli.Command{
		{
			Name:  "worker",
			Usage: "Launch Boogoo Worker",
			Action: func(c *cli.Context) error {
				db.InitDBEngine()
				if err := boogooServer.LaunchWorker(); err != nil {
					return cli.NewExitError(err.Error(), 1)
				}
				return nil
			},
		},
	}
	// Run the CLI app
	workApp.Run(os.Args)
}