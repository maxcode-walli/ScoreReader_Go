// Package p contains a Pub/Sub Cloud Function.
package p

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"cloud.google.com/go/firestore"
	firebase "firebase.google.com/go/v4"
)

var (
	client *firestore.Client
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
		log.Println("Hi Ana")
	} else {
		log.Println("Hi, person")
	}

	app, err := firebase.NewApp(ctx, &firebase.Config{ProjectID: "impactful-shard-374913"})

	if err != nil {
		log.Panic("Error init firestore")
	} else {
		client, err = app.Firestore(ctx)

		if err != nil {
			log.Panic(err)
		} else {

			err := client.RunTransaction(ctx, func(ctx context.Context, transaction *firestore.Transaction) error { //HERE TO DO

				accIter := client.Collection("users").
					Doc("K7F5Tgiucxa2NInzs3TIBe3lhyi2").
					Collection("accounts").
					Where("iban", "==", "NL02ABNA0123456789").
					Documents(ctx)

				accDoc, err := accIter.Next()
				if err != nil {
					fmt.Printf("error reading firestore account info")
				}

				person.Email = "notyo@busine.ss"

				ref := accDoc.Ref.Collection("transactions").Doc("3VhE9swaBIui8Gui14aG")

				return transaction.Create(ref, &person)
			})
			if err != nil {
				log.Panic(err, "firestore.saveTransaction")
			}

			return nil
		}
	}

	return nil
}
