package main

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"time"

	"github.com/gorilla/handlers"
	"github.com/rs/xaccess"
	"github.com/rs/xhandler"
	"github.com/rs/xlog"
	"github.com/rs/xmux"
	"github.com/rs/xstats"
	"github.com/rs/xstats/dogstatsd"
	"golang.org/x/net/context"
)

func main() {
	c := xhandler.Chain{}

	// Append a context-aware middleware handler
	c.UseC(xhandler.CloseHandler)

	// Another context-aware middleware handler
	c.UseC(xhandler.TimeoutHandler(2 * time.Second))

	// Access logs
	c.UseC(xlog.NewHandler(xlog.Config{}))

	// Stats
	flushInterval := 5 * time.Second
	tags := []string{"role:watcher"}
	statsdWriter, err := net.Dial("udp", "127.0.0.1:8126")
	if err != nil {
		log.Fatal(err)
	}
	c.UseC(xstats.NewHandler(dogstatsd.New(statsdWriter, flushInterval), tags))

	c.UseC(xaccess.NewHandler())

	c.Use(handlers.RecoveryHandler())

	mux := xmux.New()

	// Use c.Handler to terminate the chain with your final handler
	mux.GET("/welcome/:name", xhandler.HandlerFuncC(func(ctx context.Context, w http.ResponseWriter, req *http.Request) {
		fmt.Fprintf(w, "Welcome %s!", xmux.Param(ctx, "name"))
	}))

	if err := http.ListenAndServe(":8080", c.Handler(mux)); err != nil {
		log.Fatal(err)
	}
}
