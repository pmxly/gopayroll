/**
	Author: David Pan
   Date: 2020-04-21
*/

package main

import (
	"gopayroll/common"
	"gopayroll/db"
	"gopayroll/payroll/routes"
	boogooServer "gopayroll/server"
	"github.com/gin-gonic/gin"
	"github.com/urfave/cli"
	"net/http"
	"os"
	"time"
)

var (
	apiApp *cli.App
)

func init() {
	// Initialise a CLI app
	apiApp = cli.NewApp()
	apiApp.Name = "Boogoo Web Server"
	apiApp.Usage = "Process Http Request"
	apiApp.Author = "Boogoo Team"
	apiApp.Email = "tech@boogoo.cn"
	apiApp.Version = "1.0.0"
	apiApp.Flags = []cli.Flag{
		cli.StringFlag{
			Name:        "c",
			Value:       "",
			Destination: &common.ConfigPath,
			Usage:       "Path to a configuration file",
		},
	}
}

func main() {
	apiApp.Action = func(c *cli.Context) {
		db.InitDBEngine()
		_, _ = boogooServer.StartServer()
		gin.SetMode(gin.ReleaseMode)
		engine := gin.Default()
		routes.RouteManager(engine)
		s := &http.Server{
			Addr:           ":8019",
			Handler:        engine,
			ReadTimeout:    10 * time.Second,
			WriteTimeout:   10 * time.Second,
			MaxHeaderBytes: 1 << 20,
		}
		_ = s.ListenAndServe()
	}
	_ = apiApp.Run(os.Args)

	/*apiApp.Commands = []cli.Command{
		{
			Name:  "endpoint",
			Usage: "Launch Boogoo Endpoint",
			Action: func(c *cli.Context) {
				db.DB = db.InitDB()
				boogooServer.StartServer()
				engine := gin.Default()
				routes.RouteManager(engine)
				s := &http.Server{
					Addr:           ":8019",
					Handler:        engine,
					ReadTimeout:    10 * time.Second,
					WriteTimeout:   10 * time.Second,
					MaxHeaderBytes: 1 << 20,
				}
				s.ListenAndServe()
			},
		},
	}
	// Run the CLI app
	apiApp.Run(os.Args)*/
}
