package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/hibiken/asynq"
)

type Payload struct {
	URL string `json:"url"`
}

func main() {
	// Connect to Redis
	srv := asynq.NewServer(
		asynq.RedisClientOpt{Addr: "localhost:6379"},
		asynq.Config{
			Concurrency: 10,
		},
	)

	// Define task handler
	mux := asynq.NewServeMux()
	mux.HandleFunc("http:call", handlePostTask)

	// Start the worker server
	if err := srv.Run(mux); err != nil {
		panic(err)
	}
}

func handlePostTask(c context.Context, t *asynq.Task) error {
	var p Payload
	if err := json.Unmarshal(t.Payload(), &p); err != nil {
		return fmt.Errorf("failed to unmarshal payload: %w", err)
	}

	// Perform the POST request with empty body
	resp, err := http.Post(p.URL, "application/json", nil)
	if err != nil {
		return fmt.Errorf("failed to POST to %s: %w", p.URL, err)
	}
	defer resp.Body.Close()

	fmt.Println("POSTed to:", p.URL, "Status:", resp.Status)
	return nil
}
