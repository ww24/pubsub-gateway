package gateway

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"

	"cloud.google.com/go/pubsub"
	"golang.org/x/oauth2/google"
)

var (
	authorizedUsers = strings.Split(os.Getenv("AUTHORIZED_USERS"), ",")
)

type gateway struct {
	topic             *pubsub.Topic
	defaultOrigin     string
	allowOriginSuffix string
}

// New returns http.Handler.
func New(ctx context.Context) http.Handler {
	gw, err := new(ctx)
	if err != nil {
		log.Fatalln("failed to initialize gateway:", err)
	}

	mux := http.NewServeMux()

	mux.HandleFunc("/", gw.defaultHandler)
	mux.HandleFunc("/webhook", gw.publishHandler)
	mux.Handle("/publish", gw.authorizeIDToken(http.HandlerFunc(gw.publishHandler)))

	// middlewares
	handler := http.Handler(mux)
	handler = gw.cors(handler)
	return handler
}

func new(ctx context.Context) (*gateway, error) {
	cred, err := google.FindDefaultCredentials(ctx)
	if err != nil {
		return nil, fmt.Errorf("credential error: %w", err)
	}
	cli, err := pubsub.NewClient(ctx, cred.ProjectID)
	if err != nil {
		return nil, fmt.Errorf("pubsub new client error: %w", err)
	}
	return &gateway{
		topic:             cli.Topic(os.Getenv("TOPIC_NAME")),
		defaultOrigin:     os.Getenv("DEFAULT_ORIGIN"),
		allowOriginSuffix: os.Getenv("ALLOW_ORIGIN_SUFFIX"),
	}, nil
}

func (*gateway) defaultHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodOptions {
		return
	}
	w.Header().Set("Content-Type", "application/json")
	fmt.Fprintln(w, `{"status":"ok"}`)
}

func (g *gateway) publishHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodOptions:
		return
	case http.MethodPost:
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	data, err := io.ReadAll(r.Body)
	if err != nil {
		sendError(w, r, errors.New("server error"))
		log.Println("Failed to read body", err)
		return
	}

	if len(data) == 0 {
		sendBadRequestError(w, r, errors.New("payload should be specified"))
		return
	}

	ctx, cancel := context.WithCancel(r.Context())
	defer cancel()

	pr := g.topic.Publish(ctx, &pubsub.Message{Data: data})
	id, err := pr.Get(ctx)
	if err != nil {
		sendError(w, r, errors.New("server error"))
		log.Println("Failed to publish event", err)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	fmt.Fprintf(w, `{"id":"%s"}`+"\n", id)
}

func sendError(w http.ResponseWriter, r *http.Request, err error) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusInternalServerError)
	fmt.Fprintf(w, `{"error":"%s"}`+"\n", err.Error())
}

func sendBadRequestError(w http.ResponseWriter, r *http.Request, err error) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusBadRequest)
	fmt.Fprintf(w, `{"error":"%s"}`+"\n", err.Error())
}
