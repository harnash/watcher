package server

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"time"

	"github.com/gorilla/handlers"
	"github.com/harnash/watcher/logging"
	"github.com/rs/xaccess"
	"github.com/rs/xhandler"
	"github.com/rs/xlog"
	"github.com/rs/xmux"
	"github.com/rs/xstats"
	"github.com/rs/xstats/telegraf"
	"golang.org/x/net/context"
)

// Run server
func Run(serverConfig Config) {
	c := xhandler.Chain{}

	// Append a context-aware middleware handler
	c.UseC(xhandler.CloseHandler)

	// Another context-aware middleware handler
	c.UseC(xhandler.TimeoutHandler(2 * time.Second))

	// Application logs
	logging.InitConfig()
	logger := logging.GetLogger()
	c.UseC(xlog.NewHandler(logging.GetConfig()))
	log.SetFlags(0)
	log.SetOutput(logger)

	// Install some provided extra handler to set some request's context fields.
	// Thanks to those handler, all our logs will come with some pre-populated fields.
	c.UseC(xlog.MethodHandler("method"))
	c.UseC(xlog.URLHandler("url"))
	c.UseC(xlog.RemoteAddrHandler("ip"))
	c.UseC(xlog.UserAgentHandler("user_agent"))
	c.UseC(xlog.RefererHandler("referer"))
	c.UseC(xlog.RequestIDHandler("req_id", "Request-Id"))

	// Stats
	flushInterval := 5 * time.Second
	tags := []string{"role:watcher"}
	telegrafWriter, err := net.Dial("udp", "127.0.0.1:8125")
	if err != nil {
		log.Fatal(err)
	}
	c.UseC(xstats.NewHandler(telegraf.New(telegrafWriter, flushInterval), tags))

	c.UseC(xaccess.NewHandler())

	c.Use(handlers.RecoveryHandler())

	mux := xmux.New()

	// Use c.Handler to terminate the chain with your final handler
	mux.GET("/welcome/:name", xhandler.HandlerFuncC(func(ctx context.Context, w http.ResponseWriter, req *http.Request) {
		userID := xmux.Param(ctx, "name")

		// Get the logger from the context. You can safely assume it will be always there,
		// if the handler is removed, xlog.FromContext will return a NopLogger
		l := xlog.FromContext(ctx)

		l.Info("User welcomed", xlog.F{
			"user":   userID,
			"status": "ok",
		})

		fmt.Fprintf(w, "Welcome %s!", userID)
	}))

	logger.Infof("Listening on: %s", serverConfig.ListenAddress)

	if err := http.ListenAndServe(serverConfig.ListenAddress, c.Handler(mux)); err != nil {
		log.Fatal(err)
	}
}
