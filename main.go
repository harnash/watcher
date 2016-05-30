package main

import (
	"os"

	"github.com/harnash/watcher/server"
	"github.com/jawher/mow.cli"
)

func main() {
	app := cli.App("watcher", "Run server")

	app.Spec = "[-l]"

	var (
		listenAddr = app.StringOpt("l listen", ":8080", "address to listen on")
	)
	app.Action = func() {
		var conf server.Config

		conf.ListenAddress = *listenAddr

		server.Run(conf)
	}

	app.Run(os.Args)
}
