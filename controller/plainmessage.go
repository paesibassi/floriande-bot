package controller

import (
	"fmt"
	"log"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"gitlab.com/gruppi-preparazione/floriande-bot/store"
)

func handleBook(update *tgbotapi.Update, eventID int) tgbotapi.MessageConfig {
	store.AddBooking(client, eventID, update.Message.From.ID, update.Message.From.FirstName)
	confirm := fmt.Sprintf("%v %s just booked for [%d]", bell, update.Message.From.FirstName, eventID)
	if _, err := bot.Send(tgbotapi.NewMessage(barmanID, confirm)); err != nil {
		log.Fatal(err)
	}
	return tgbotapi.NewMessage(update.Message.Chat.ID, fmt.Sprintf("Great! We reserved a spot for you %s!", update.Message.From.FirstName))
}

func handleEcho(update *tgbotapi.Update) tgbotapi.MessageConfig {
	// simple echo
	return tgbotapi.NewMessage(update.Message.Chat.ID, fmt.Sprintf("%v %s", wineGlass, update.Message.Text))
}

func handlePlainMessage(update *tgbotapi.Update) tgbotapi.MessageConfig {
	if id, ok := store.Events[update.Message.Text]; ok {
		return handleBook(update, id)
	} else {
		return handleEcho(update)
	}
}
