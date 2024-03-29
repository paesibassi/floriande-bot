package store

import (
	"context"
	"fmt"
	"log"
	"time"

	"cloud.google.com/go/firestore"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

var Events = map[string]int{
	"Federico Birthday May 5h":       20230505,
	"Spring Cocktail Party May 19th": 20230519,
}

func NewEventsKeyboard(events map[string]int) tgbotapi.ReplyKeyboardMarkup {
	rows := make([][]tgbotapi.KeyboardButton, 0, len(events))
	for event := range events {
		rows = append(rows, tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton(event),
		))
	}
	return tgbotapi.NewOneTimeReplyKeyboard(rows...)
}

var EventsKeyboard = NewEventsKeyboard(Events)

type Booking struct {
	BookingID string
	Timestamp time.Time
	Customer
	EventID int
}

func NewBooking(eventID int, customerID int64, customerName, customerLang string) Booking {
	t := time.Now()
	bookingID := fmt.Sprintf("#%d%s", eventID, customerName)
	return Booking{
		BookingID: bookingID,
		Timestamp: t,
		Customer: Customer{
			CustomerID:   customerID,
			CustomerName: customerName,
			CustomerLang: customerLang,
		},
		EventID: eventID,
	}
}

func AddBooking(client *firestore.Client, eventID int, customerID int64, customerName, customerLang string) {
	b := NewBooking(eventID, customerID, customerName, customerLang)
	ctx := context.Background()
	_, err := client.
		Collection("bookings").
		Doc(b.BookingID).
		Set(ctx, b)
	if err != nil {
		log.Fatalf("Failed adding booking: %v", err)
	}
}
