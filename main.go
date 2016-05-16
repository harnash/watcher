package main

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/handlers"
	"github.com/rs/xaccess"
	"github.com/rs/xhandler"
	"github.com/rs/xlog"
	"github.com/rs/xmux"
	"github.com/rs/xstats"
	"github.com/rs/xstats/telegraf"
	"golang.org/x/net/context"
)

func main() {
	c := xhandler.Chain{}

	// Append a context-aware middleware handler
	c.UseC(xhandler.CloseHandler)

	// Another context-aware middleware handler
	c.UseC(xhandler.TimeoutHandler(2 * time.Second))

	host, _ := os.Hostname()
	// logstashWriter, err := net.Dial("udp", "127.0.0.1:1410")
	// if err != nil {
	// 	log.Fatal(err)
	// }
	conf := xlog.Config{
		// Log info level and higher
		Level: xlog.LevelInfo,
		// Set some global env fields
		Fields: xlog.F{
			"role": "watcher",
			"host": host,
		},
		// Output everything on console
		Output: xlog.NewOutputChannel(xlog.MultiOutput{
			0: xlog.NewConsoleOutput(),
			// 1: xlog.NewLogstashOutput(logstashWriter),
		}),
	}

	// Application logs
	log.SetFlags(0)
	xlogger := xlog.New(conf)
	c.UseC(xlog.NewHandler(conf))
	log.SetOutput(xlogger)

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

	if err := http.ListenAndServe(":8080", c.Handler(mux)); err != nil {
		log.Fatal(err)
	}
}
