package main

import (
	"SmartLocker/cmd/server/router"
	"SmartLocker/config"
	"SmartLocker/database"
	"SmartLocker/logger"
	"github.com/go-playground/log"
	"github.com/urfave/cli"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	// init the helpers
	logger.Setup()
	config.Setup()
	database.Setup()

	// create an app instance
	app := cli.NewApp()
	app.Name = "SmartLocker"
	app.Usage = "the backend of SmartLocker"
	app.Version = getVersion()
	app.Commands = []cli.Command{
		{
			Name:    "server",
			Aliases: []string{"s"},
			Usage:   "Start the Server",
			Action:  StartServer,
		},
	}

	err := app.Run(os.Args)

	if err != nil {
		log.WithError(err).Fatal("Couldn't startup")
	}

}

func StartServer(ctx *cli.Context) {
	// init the route
	routersInit := router.InitRouters()
	addr := config.Conf.WebServer.Address + ":" + config.Conf.WebServer.Port

	// setup a http server instance
	server := &http.Server{
		Addr:    addr,
		Handler: routersInit,
	}

	// start the server
	go func() {
		err := server.ListenAndServe()
		if err != nil {
			log.WithError(err).Fatal("Couldn't start server")
		}
	}()

	// shutdown
	// Handle SIGINT and SIGTERM.
	ch := make(chan os.Signal)
	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)
	log.Debug(<-ch)

	// Stop the service gracefully.
	database.CloseDB()
	err := server.Shutdown(nil)
	if err != nil {
		log.WithError(err).Warn("Couldn't shutdown the server")
	}
}

var (
	BuildTimeStamp = "None"
	GitHash        = "None"
	Version        = "0.1"
)

func getVersion() string {
	return Version + "+" + GitHash + "+" + BuildTimeStamp
}
