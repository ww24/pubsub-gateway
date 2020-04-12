// +build cloudrun

package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/ww24/pubsub-gateway/gateway"
	"github.com/ww24/pubsub-gateway/receiver"
)

const (
	defaultPort     = "8080"
	defaultMode     = "gateway"
	shutdownTimeout = 30 * time.Second
)

var (
	port     = os.Getenv("PORT")
	mode     = os.Getenv("MODE")
	confFile = flag.String("config", "", "set path to config (required for receiver mode)")
)

func main() {
	flag.Parse()
	if port == "" {
		port = defaultPort
	}
	if mode == "" {
		mode = defaultMode
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	var closeFn func(context.Context)
	if mode == "" {
		mode = "gateway"
	}
	switch mode {
	case "gateway":
		closeFn = gatewayMode(ctx, port)

	case "receiver":
		closeFn = receiverMode(ctx)
	}

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM, syscall.SIGHUP)
	<-sigCh

	ctx, cancel = context.WithTimeout(context.Background(), shutdownTimeout)
	defer cancel()
	closeFn(ctx)
}

func gatewayMode(ctx context.Context, port string) func(context.Context) {
	ctx, cancel := context.WithCancel(ctx)
	log.Println("Listen at", port)
	srv := &http.Server{
		Addr:    ":" + port,
		Handler: gateway.New(ctx),
	}
	go func() {
		if err := srv.ListenAndServe(); err != nil {
			log.Println("Server Error:", err)
		}
	}()
	return func(ctx context.Context) {
		cancel()
		if err := srv.Shutdown(ctx); err != nil {
			log.Println("Shutdown Error:", err)
		}
	}
}

func receiverMode(ctx context.Context) func(context.Context) {
	if *confFile == "" {
		fmt.Fprintln(os.Stderr, "-config flag is required")
		fmt.Fprintf(os.Stderr, "Usage of %s:\n", os.Args[0])
		flag.PrintDefaults()
		os.Exit(1)
	}

	ctx, cancel := context.WithCancel(ctx)
	config, err := receiver.Parse(*confFile)
	if err != nil {
		cancel()
		log.Fatalln("failed to load config:", err)
		return func(context.Context) {}
	}
	s, err := receiver.New(ctx)
	if err != nil {
		cancel()
		log.Fatalln("failed initialize receiver:", err)
		return func(context.Context) {}
	}
	for _, handler := range config.Handlers {
		fmt.Println("subscription:", handler.Subscription)
		handler := handler
		go func() {
			err := s.Subscribe(ctx, handler.Subscription, func(ctx context.Context, data []byte) bool {
				log.Println("received:", string(data))
				m := make(map[string]interface{})
				if err := json.Unmarshal(data, &m); err != nil {
					log.Println("failed to unmarshal json:", err)
					return true
				}

				var action receiver.Executable
				var payload []byte
				switch handler.Action.Type {
				case receiver.ActionHTTP:
					a := handler.Action.HTTPRequestAction
					action = receiver.NewHTTPAction(a.Header, a.Method, a.URL)
					var err error
					payload, err = a.Payload.RenderJSON(data)
					if err != nil {
						log.Println("failed to render:", err)
						return false
					}

				default:
					log.Println("unexpected action type:", handler.Action.Type)
					return false
				}

				log.Println("payload:", string(payload))
				if err := action.Exec(ctx, payload); err != nil {
					log.Println("failed to exec action", err)
					return false
				}

				return true
			})
			if err != nil {
				panic(err)
			}
		}()
	}
	return func(context.Context) {
		cancel()
	}
}
