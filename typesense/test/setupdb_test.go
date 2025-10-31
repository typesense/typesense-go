//go:build integration && !docker
// +build integration,!docker

package test

import (
	"context"
	"errors"
	"os"
	"time"

	"github.com/typesense/typesense-go/v4/typesense"
)

func waitHealthyStatus(client *typesense.Client, timeout time.Duration) error {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(1 * time.Second):
			if healthy, _ := client.Health(context.Background(), 2*time.Second); !healthy {
				continue
			}
			return nil
		}
	}
}

func setupDB() (*typesense.Client, error) {
	url := os.Getenv("TYPESENSE_URL")
	apiKey := os.Getenv("TYPESENSE_API_KEY")
	if len(url) == 0 || len(apiKey) == 0 {
		return nil, errors.New("TYPESENSE_URL or TYPESENSE_API_KEY env variable is empty!")
	}
	client := typesense.NewClient(
		typesense.WithServer(url),
		typesense.WithAPIKey(apiKey))
	if err := waitHealthyStatus(client, 1*time.Minute); err != nil {
		return nil, err
	}
	return client, nil
}

func stopDB() error {
	return nil
}
