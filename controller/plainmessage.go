package controller

import (
	"fmt"
	"log"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"gitlab.com/gruppi-preparazione/floriande-bot/store"
)

type PlainMsgInGroupError struct {
	message string
}

func newPlainMsgInGroup(message string) *PlainMsgInGroupError {
	return &PlainMsgInGroupError{
		message: message,
	}
}

func (e *PlainMsgInGroupError) Error() string {
	return e.message
}

func handlePlainMessage(update *tgbotapi.Update) error {
	var msg tgbotapi.MessageConfig
	if update.Message.Chat.Type != "private" {
		return newPlainMsgInGroup("ignore plain messages sent in groups")
	}
	if id, ok := store.Events[update.Message.Text]; ok {
		msg = handleBook(update, id)
	} else if category, ok := store.AllCocktails[update.Message.Text]; ok {
		msg = handleOrderPlainMessage(update, category)
	} else {
		msg = handleDidntUnderstand(update)
	}
	_, err := bot.Send(msg)
	return err
}

func handleOrderPlainMessage(update *tgbotapi.Update, category string) tgbotapi.MessageConfig {
	drink := update.Message.Text
	store.AddOrder(client, update.Message.From.ID, update.Message.From.FirstName, drink, category)
	confirm := fmt.Sprintf("%v %s just ordered a %s", bell, update.Message.From.FirstName, drink)
	if _, err := bot.Send(tgbotapi.NewMessage(barmanID, confirm)); err != nil {
		log.Fatal(err)
	}
	return tgbotapi.NewMessage(
		update.Message.From.ID,
		fmt.Sprintf("A %s is coming soon! %v", drink, cocktailGlass),
	)
}

func handleBook(update *tgbotapi.Update, eventID int) tgbotapi.MessageConfig {
	store.AddBooking(client, eventID, update.Message.From.ID, update.Message.From.FirstName)
	confirm := fmt.Sprintf("%v %s just booked for [%d]", bell, update.Message.From.FirstName, eventID)
	if _, err := bot.Send(tgbotapi.NewMessage(barmanID, confirm)); err != nil {
		log.Fatal(err)
	}
	return tgbotapi.NewMessage(update.Message.Chat.ID, fmt.Sprintf("Great! We reserved a spot for you %s!", update.Message.From.FirstName))
}

var didntUnderstandMsg = `I didn't understand: '%s'\.
Please try the guided order command ` + "`/drink`" + ` to get help while ordering\.
Otherwise, make sure to type the exact name of the drink\.
Get in touch with [us](tg://user?id=%d) if anything doesn't work\!`

func handleDidntUnderstand(update *tgbotapi.Update) tgbotapi.MessageConfig {
	msg := tgbotapi.NewMessage(update.Message.Chat.ID, fmt.Sprintf(didntUnderstandMsg, update.Message.Text, barmanID))
	msg.ParseMode = "MarkdownV2"
	return msg
}
