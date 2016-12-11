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
	"io/ioutil"
	"log"
	"net/http"

	"cloud.google.com/go/pubsub"
	"cloud.google.com/go/trace"
)

// Publish request body to a Google Cloud Pub/Sub topic.
type pubSubHandler struct {
	topic        string
	pubsubClient *pubsub.Client
	traceClient  *trace.Client
}

// PubSubHandler returns a request handler that publishes
// each request body it receives to the given pub/sub topic.
func PubSubHandler(topic string, pubsubClient *pubsub.Client, traceClient *trace.Client) http.Handler {
	return &pubSubHandler{topic, pubsubClient, traceClient}
}

func (ph *pubSubHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	span := ph.traceClient.SpanFromRequest(r)
	defer span.Finish()

	topic := ph.pubsubClient.Topic(ph.topic)

	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Println("Failed to extract message from request: %v", err)
		http.Error(w, "Failed to extract message from request", http.StatusInternalServerError)
		return
	}

	childSpan := span.NewChild("pubsub")
	childSpan.SetLabel("topic", ph.topic)

	ctx := context.Background()
	msgIDs, err := topic.Publish(ctx, &pubsub.Message{
		Data: data,
	})

	childSpan.Finish()

	if err != nil {
		log.Printf("Failed to publish message: %v", err)
		http.Error(w, "Failed to publish message", http.StatusInternalServerError)
		return
	}

	log.Printf("Published a message with a message ID: %s\n", msgIDs[0])
}
