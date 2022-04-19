package store

import (
	"context"
	"log"

	"cloud.google.com/go/firestore"
	firebase "firebase.google.com/go"
)

func SetupFirestore() *firestore.Client {
	// Use the application default credentials
	ctx := context.Background()
	conf := &firebase.Config{ProjectID: "floriande-bot"}
	app, err := firebase.NewApp(ctx, conf)
	if err != nil {
		log.Fatalln(err)
	}

	client, err := app.Firestore(ctx)
	if err != nil {
		log.Fatalln(err)
	}
	return client
}
