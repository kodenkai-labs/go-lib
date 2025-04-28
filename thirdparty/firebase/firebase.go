package firebase

import (
	"context"
	"fmt"

	"cloud.google.com/go/firestore"
	firebase "firebase.google.com/go"
	"google.golang.org/api/option"
)

type Client struct {
	client *firestore.Client
}

func NewClient(ctx context.Context, firebaseCfgPath string) (*Client, error) {
	app, err := firebase.NewApp(ctx, nil, option.WithCredentialsFile(firebaseCfgPath))
	if err != nil {
		return nil, fmt.Errorf("new firebase app: %w", err)
	}

	client, err := app.Firestore(ctx)
	if err != nil {
		return nil, fmt.Errorf("new firestore client: %w", err)
	}

	return &Client{
		client: client,
	}, nil
}

func (c *Client) Close() error {
	return c.client.Close()
}

func (c *Client) AddDocument(ctx context.Context, collection string, data map[string]any) error {
	_, _, err := c.client.Collection(collection).Add(ctx, data)
	if err != nil {
		return fmt.Errorf("add document: %w", err)
	}

	return nil
}
