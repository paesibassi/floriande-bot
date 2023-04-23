package floriandebot

import (
	"fmt"
	"log"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"gitlab.com/gruppi-preparazione/floriande-bot/store"
)

type PlainMsgInGroupError struct {
	message string
}

func newPlainMsgInGroupError(message string) *PlainMsgInGroupError {
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
		return newPlainMsgInGroupError("ignore plain messages sent in groups")
	}

	allCocktails := store.NewMenuList(store.CocktailsMenu(client))
	if id, ok := store.Events[update.Message.Text]; ok {
		msg = handleBook(update, id)
	} else if category, ok := allCocktails[update.Message.Text]; ok {
		msg = handleOrderPlainMessage(update, category)
	} else if update.Message.Photo != nil {
		msg = handleReturnFileId(update)
	} else {
		msg = handleDontUnderstand(update)
	}
	_, err := bot.Send(msg)
	return err
}

func handleOrderPlainMessage(update *tgbotapi.Update, category string) tgbotapi.MessageConfig {
	drink := update.Message.Text
	store.AddOrder(client,
		update.Message.From.ID,
		update.Message.From.FirstName,
		update.Message.From.LanguageCode, drink,
		category,
	)
	confirm := fmt.Sprintf("%v %s just ordered a %s", bell, update.Message.From.FirstName, drink)
	if _, err := bot.Send(tgbotapi.NewMessage(barmanID, confirm)); err != nil {
		log.Fatal(err)
	}
	msgText := mss[orderConfirmation][userLanguage(update.Message.From.LanguageCode)]
	return tgbotapi.NewMessage(
		update.Message.From.ID,
		fmt.Sprintf(msgText, drink, cocktailGlass),
	)
}

func handleBook(update *tgbotapi.Update, eventID int) tgbotapi.MessageConfig {
	store.AddBooking(client, eventID, update.Message.From.ID, update.Message.From.FirstName, update.Message.From.LanguageCode)
	confirm := fmt.Sprintf("%v %s just booked for [%d]", bell, update.Message.From.FirstName, eventID)
	if _, err := bot.Send(tgbotapi.NewMessage(barmanID, confirm)); err != nil {
		log.Fatal(err)
	}
	msgText := mss[bookConfirmation][userLanguage(update.Message.From.LanguageCode)]
	return tgbotapi.NewMessage(update.Message.Chat.ID, fmt.Sprintf(msgText, eventID, update.Message.From.FirstName))
}

func handleReturnFileId(update *tgbotapi.Update) tgbotapi.MessageConfig {
	msgText := mss[returnFileId][userLanguage(update.Message.From.LanguageCode)]
	var fileId string
	var lastImageWidth int
	for _, f := range update.Message.Photo {
		if f.Width > lastImageWidth {
			fileId = f.FileID
			lastImageWidth = f.Width
		}
	}
	msgText = fmt.Sprintf(msgText, fileId)
	if update.Message.Caption != "" {
		msgText += fmt.Sprintf(": '%s'", update.Message.Caption)
	}
	msg := tgbotapi.NewMessage(update.Message.Chat.ID, msgText)
	return msg
}

func handleDontUnderstand(update *tgbotapi.Update) tgbotapi.MessageConfig {
	msgText := mss[dontUnderstand][userLanguage(update.Message.From.LanguageCode)]
	msg := tgbotapi.NewMessage(update.Message.Chat.ID, fmt.Sprintf(msgText, barmanID))
	msg.ParseMode = "MarkdownV2"
	return msg
}
