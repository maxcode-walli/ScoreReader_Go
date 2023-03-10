// Package p contains a Pub/Sub Cloud Function.
package p

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strings"

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
	var transactionScore TransactionScore

	//event := m.ToPigeonEvent()
	_ = json.Unmarshal(m.Data, &transactionScore)

	app, err := firebase.NewApp(ctx, &firebase.Config{ProjectID: "impactful-shard-374913"})

	if err != nil {
		log.Panic("Error init firestore")
	} else {
		client, err = app.Firestore(ctx)

		if err != nil {
			log.Panic(err)
		} else {

			err := client.RunTransaction(ctx, func(ctx context.Context, transaction *firestore.Transaction) error {

				accIter := client.Collection("users").Doc(transactionScore.UserID).
					Collection("accounts").
					Where("externalAccountId", "==", transactionScore.ExternalAccountId).
					Documents(ctx)

				accDoc, err := accIter.Next()
				if err != nil {
					log.Panic(fmt.Sprintf("no user found with externalAccountId: %s", transactionScore.ExternalAccountId))
				}
				ref := accDoc.Ref.Collection("transactions").Doc(transactionScore.TransactionID)

				if strings.Contains(transactionScore.Label, "risk") {
					var updates []firestore.Update
					updates = append(updates, firestore.Update{
						Path:  "IsAnomaly",
						Value: true,
					})
					transaction.Update(ref, updates)
				} else if strings.Contains(transactionScore.Label, "legit") {
					var updates []firestore.Update
					updates = append(updates, firestore.Update{
						Path:  "IsAnomaly",
						Value: false,
					})
					transaction.Update(ref, updates)
				}

				return nil
			})
			if err != nil {
				log.Panic(err, "firestore.saveTransaction")
			}

			return nil
		}
	}

	return nil
}
