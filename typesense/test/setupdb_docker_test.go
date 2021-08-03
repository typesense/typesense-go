// +build integration,docker

package test

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"time"

	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
	"github.com/typesense/typesense-go/typesense"
	"github.com/typesense/typesense-go/typesense/api"
)

var typesenseC testcontainers.Container

func setupDB() (client *typesense.Client, err error) {
	log.Println("starting typesense container...")
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Minute)
	defer cancel()

	testAPIKey := "abcdef"
	req := testcontainers.ContainerRequest{
		Image:        "typesense/typesense:0.21.0",
		ExposedPorts: []string{"8108/tcp"},
		AutoRemove:   false,
		Cmd: []string{"--data-dir", "/tmp",
			fmt.Sprintf("--api-key=%s", testAPIKey)},
		WaitingFor: wait.ForHTTP("/health").WithPort("8108/tcp").
			WithStartupTimeout(1 * time.Minute).
			WithPollInterval(1 * time.Second).
			WithResponseMatcher(func(body io.Reader) bool {
				jd := json.NewDecoder(body)
				var result api.HealthStatus
				if err := jd.Decode(&result); err != nil {
					return false
				}
				return result.Ok
			}),
	}
	typesenseC, err = testcontainers.GenericContainer(ctx,
		testcontainers.GenericContainerRequest{
			ContainerRequest: req,
			Started:          true,
		})
	if err != nil {
		log.Println("container error!")
		return nil, err
	}
	dbURL, err := typesenseC.PortEndpoint(ctx, "8108/tcp", "http")
	if err != nil {
		return nil, err
	}

	client = typesense.NewClient(
		typesense.WithServer(dbURL),
		typesense.WithAPIKey(testAPIKey),
	)
	return client, nil
}

func stopDB() error {
	log.Println("terminating typesense container...")
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Minute)
	defer cancel()
	return typesenseC.Terminate(ctx)
}
