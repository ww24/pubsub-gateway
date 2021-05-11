package receiver

import (
	"context"
	"fmt"
	"os"

	"cloud.google.com/go/pubsub"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/option"
)

// Subscriber provides pubsub subscriber.
type Subscriber interface {
	Subscribe(ctx context.Context, subscription string, fn func(context.Context, []byte) bool) error
}

type receiver struct {
	cli *pubsub.Client
}

// New returns Subscriber.
func New(ctx context.Context) (Subscriber, error) {
	var cred *google.Credentials
	var err error
	if data := os.Getenv("SERVICE_ACCOUNT_JSON"); data != "" {
		cred, err = google.CredentialsFromJSON(ctx, []byte(data), pubsub.ScopePubSub)
	} else {
		cred, err = google.FindDefaultCredentials(ctx, pubsub.ScopePubSub)
	}
	if err != nil {
		return nil, fmt.Errorf("credential error: %w", err)
	}

	cli, err := pubsub.NewClient(ctx, cred.ProjectID, option.WithCredentials(cred))
	if err != nil {
		return nil, fmt.Errorf("pubsub new client error: %w", err)
	}

	return &receiver{
		cli: cli,
	}, nil
}

func (r *receiver) Subscribe(ctx context.Context, subscription string, fn func(context.Context, []byte) bool) error {
	return r.cli.Subscription(subscription).Receive(ctx, func(ctx context.Context, m *pubsub.Message) {
		if fn(ctx, m.Data) {
			m.Ack()
		}
	})
}
