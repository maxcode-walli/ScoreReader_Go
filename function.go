// Package p contains a Pub/Sub Cloud Function.
package p

import (
	"context"
	"encoding/json"
	"log"
)

// PubSubMessage is the payload of a Pub/Sub event. Please refer to the docs for
// additional information regarding Pub/Sub events.
type PubSubMessage struct {
	Data []byte `json:"data"`
}

// HelloPubSub consumes a Pub/Sub message.
func HelloPubSub(ctx context.Context, m PubSubMessage) error {
	log.Println(string(m.Data))
	var person Person
	//event := m.ToPigeonEvent()
	_ = json.Unmarshal(m.Data, &person)

	if person.Name == "Ana" {
		log.Println("Hi Alex")
	} else {
		log.Println("Hi, person")
	}
	return nil
}
