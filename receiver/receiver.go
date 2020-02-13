package receiver

import (
	"context"
	"fmt"
	"os"

	"cloud.google.com/go/pubsub"
	"golang.org/x/oauth2/google"
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
	cred, err := google.FindDefaultCredentials(ctx)
	if err != nil {
		return nil, fmt.Errorf("credential error: %w", err)
	}
	projectID := cred.ProjectID
	if projectID == "" {
		projectID = os.Getenv("PROJECT_ID")
	}
	cli, err := pubsub.NewClient(ctx, projectID)
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
