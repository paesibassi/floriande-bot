package main

import (
	"context"
	"encoding/json"
	"flag"
	"io/ioutil"
	"log"
	"os"

	"cloud.google.com/go/firestore"
	firebase "firebase.google.com/go"
	"gitlab.com/gruppi-preparazione/floriande-bot/store"
	"google.golang.org/api/iterator"
	"google.golang.org/api/option"
)

type Menu struct {
	Menu []Category `json:"menu"`
}

type Category struct {
	Name      string     `json:"name"`
	Cocktails []Cocktail `json:"cocktails"`
}

type Cocktail struct {
	Name  string `json:"name"`
	Image string `json:"image,omitempty"`
}

func main() {
	flag.Parse()
	args := flag.Args()

	menu, err := readMenuFromFile(args[0])
	if err != nil {
		log.Println(err)
		return
	}

	client := setupFirestoreClient()
	defer client.Close()
	err = storeMenuInFirestore(client, menu)
	if err != nil {
		log.Println(err)
	}
}

func readMenuFromFile(filename string) (*Menu, error) {
	log.Printf("Open menu data from '%s'\n", filename)

	f, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	bytes, _ := ioutil.ReadAll(f)

	var menu Menu
	err = json.Unmarshal(bytes, &menu)
	if err != nil {
		return nil, err
	}

	return &menu, nil
}

func setupFirestoreClient() *firestore.Client {
	ctx := context.Background()
	conf := &firebase.Config{ProjectID: "floriande-bot"}
	// Firebase firestore service account token
	opt := option.WithCredentialsFile("./service/floriande-bot-firebase-adminsdk.json")
	app, err := firebase.NewApp(ctx, conf, opt)
	if err != nil {
		log.Fatalln(err)
	}

	client, err := app.Firestore(ctx)
	if err != nil {
		log.Fatalln(err)
	}
	return client
}

func storeMenuInFirestore(client *firestore.Client, menu *Menu) error {
	ctx := context.Background()

	nCategories := len(menu.Menu)
	var nCocktails int
	for _, cat := range menu.Menu {
		nCocktails += len(cat.Cocktails)
	}
	log.Printf("Updating %d cocktails in %d categories\n", nCocktails, nCategories)

	coll := client.Collection("menu")
	batch := client.Batch()

	// Delete all existing items before committing the new ones
	iter := coll.Documents(context.Background())
	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			log.Fatalf("Failed to iterate: %v", err)
		}
		batch.Delete(doc.Ref)
	}

	// Add the new menu items
	for _, category := range menu.Menu {
		for _, cocktail := range category.Cocktails {
			docRef := coll.Doc(cocktail.Name)
			c := store.NewCocktail(cocktail.Name, category.Name)
			batch.Set(docRef, c)
		}
	}

	// Commit all the menu changes as a single batch operation
	_, err := batch.Commit(ctx)
	return err
}
