package store

import (
	"bytes"
	"context"
	"fmt"
	"log"
	"text/tabwriter"
	"time"

	"cloud.google.com/go/firestore"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"google.golang.org/api/iterator"
)

type Customer struct {
	CustomerID   int64
	CustomerName string
	CustomerLang string
}

type Order struct {
	OrderID   string
	Timestamp time.Time
	Customer
	Cocktail
	Quantity int
	Done     bool
}

func NewOrder(customerID int64, customerName, customerLang, drink, category string) Order {
	t := time.Now()
	orderID := fmt.Sprintf("#%d%s", t.Unix(), customerName)
	cocktail := NewCocktail(drink, category)
	return Order{
		OrderID:   orderID,
		Timestamp: t,
		Customer: Customer{
			CustomerID:   customerID,
			CustomerName: customerName,
			CustomerLang: customerLang,
		},
		Cocktail: cocktail,
		Quantity: 1,
		Done:     false,
	}
}

func AddOrder(client *firestore.Client, customerID int64, customerName, customerLang, drink, category string) {
	o := NewOrder(customerID, customerName, customerLang, drink, category)
	ctx := context.Background()
	_, err := client.
		Collection("orders").
		Doc(o.OrderID).
		Set(ctx, o)
	if err != nil {
		log.Fatalf("Failed adding order: %v", err)
	}
}

func OrderDetails(client *firestore.Client, orderID string) (*Order, error) {
	orderPtr := new(Order)
	dsnap, err := client.
		Collection("orders").
		Doc(orderID).
		Get(context.Background())
	if err != nil {
		return nil, err
	}
	err = dsnap.DataTo(orderPtr)
	if err != nil {
		return nil, err
	}
	return orderPtr, err
}

func CloseOrder(client *firestore.Client, orderID string) error {
	_, err := client.
		Collection("orders").
		Doc(orderID).
		Update(context.Background(), []firestore.Update{
			{
				Path:  "Done",
				Value: true,
			},
		})
	return err
}

type Orders []Order

func (o Orders) String() string {
	var b bytes.Buffer
	w := tabwriter.NewWriter(&b, 0, 0, 2, ' ', 0)
	for i, o := range o {
		fmt.Fprintf(w, "%1.1d\t%0.12s\t for %0.9s\t@ %s\n", i+1, o.Cocktail.CocktailName, o.CustomerName, o.Timestamp.Format("15:04"))
	}
	w.Flush()
	return b.String()
}

func UserOrders(client *firestore.Client, id int64) Orders {
	var order Order
	var orders Orders
	iter := client.
		Collection("orders").
		Where("CustomerID", "==", id).
		Where("Done", "==", false).
		OrderBy("Timestamp", firestore.Asc).
		Documents(context.Background())
	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			log.Fatalf("Failed to iterate: %v", err)
		}
		err = doc.DataTo(&order)
		if err != nil {
			continue
		}
		orders = append(orders, order)
	}
	fmt.Printf("%v", orders)
	return orders
}

func AllOrders(client *firestore.Client) Orders {
	var order Order
	var orders Orders
	iter := client.
		Collection("orders").
		Where("Done", "==", false).
		OrderBy("Timestamp", firestore.Asc).
		Documents(context.Background())
	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			log.Fatalf("Failed to iterate: %v", err)
		}
		err = doc.DataTo(&order)
		if err != nil {
			continue
		}
		orders = append(orders, order)
	}
	return orders
}

func NewServeOrdersKeyboard(orders Orders) tgbotapi.InlineKeyboardMarkup {
	const buttonsPerRow = 3

	nButtons := len(orders)
	buttons := make([]tgbotapi.InlineKeyboardButton, nButtons)
	for i, o := range orders {
		text := fmt.Sprintf("%2.1d %3.3s > %5.3s", i, o.CocktailName, o.CustomerName)
		buttons[i] = tgbotapi.NewInlineKeyboardButtonData(text, o.OrderID)
	}
	var rows [][]tgbotapi.InlineKeyboardButton
	for i := 0; i < nButtons; i += buttonsPerRow {
		limit := i + buttonsPerRow
		if limit > nButtons {
			limit = nButtons
		}
		rows = append(rows, tgbotapi.NewInlineKeyboardRow(buttons[i:limit]...))
	}
	return tgbotapi.NewInlineKeyboardMarkup(rows...)
}
