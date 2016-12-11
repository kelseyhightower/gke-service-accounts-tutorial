// Copyright 2016 Google Inc. All Rights Reserved.
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//     http://www.apache.org/licenses/LICENSE-2.0
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"cloud.google.com/go/pubsub"
	"cloud.google.com/go/trace"
)

func main() {
	projectID := os.Getenv("PROJECT_ID")
	if projectID == "" {
		log.Fatalf("PROJECT_ID must be set and non-empty.")
	}

	if os.Getenv("GOOGLE_APPLICATION_CREDENTIALS") == "" {
		log.Fatalf("GOOGLE_APPLICATION_CREDENTIALS must be set and non-empty.")
	}

	topic := os.Getenv("TOPIC")
	if topic == "" {
		log.Fatalf("TOPIC must be set and non-empty")
	}

	tctx := context.Background()
	traceClient, err := trace.NewClient(tctx, projectID)
	if err != nil {
		log.Fatalf("Failed to create trace client: %v", err)
	}

	p, err := trace.NewLimitedSampler(0.1, 5)
	if err != nil {
		log.Fatalf("Failed to set tracing sampling policy: %v", err)
	}
	traceClient.SetSamplingPolicy(p)

	pctx := context.Background()
	pubsubClient, err := pubsub.NewClient(pctx, projectID)
	if err != nil {
		log.Fatalf("Failed to create pubsub client: %v", err)
	}

	http.Handle("/pubsub", PubSubHandler(topic, pubsubClient, traceClient))

	server := &http.Server{Addr: ":8080"}
	go func() {
		log.Fatal(server.ListenAndServe())
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Printf("Shutdown signal received, exiting...")
	server.Shutdown(context.Background())
}
