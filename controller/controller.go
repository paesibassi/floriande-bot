package controller

import (
	"fmt"
	"log"

	"cloud.google.com/go/firestore"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

var client *firestore.Client
var bot *tgbotapi.BotAPI

const (
	bell            = "\xF0\x9F\x94\x94"
	clinkingGlasses = "\xF0\x9F\xA5\x82"
	cocktailGlass   = "\xF0\x9F\x8D\xB8"
	personTipping   = "\xF0\x9F\x92\x81"
	rainbow         = "\xF0\x9F\x8C\x88"
	sun             = "\xE2\x98\x80"
	tropicalDrink   = "\xF0\x9F\x8D\xB9"
	tumblerGlass    = "\xF0\x9F\xA5\x83"
	wineGlass       = "\xF0\x9F\x8D\xB7"

	barmanID = 20137373 // admin chat_id
)

func Setup(b *tgbotapi.BotAPI, cl *firestore.Client) {
	bot = b
	client = cl
}

func CloseDB() {
	client.Close()
}

func HandleUpdate(update *tgbotapi.Update) error {
	var err error
	switch {
	case update.Message != nil:
		var msg tgbotapi.MessageConfig
		if update.Message.IsCommand() {
			msg = handleCommands(update)
		} else {
			msg = handlePlainMessage(update)
		}
		_, err = bot.Send(msg)
	case update.CallbackQuery != nil:
		handleCallBackQuery(update)
	case update.MyChatMember != nil:
		msg := handleChatMemberUpdate(update)
		_, err = bot.Send(msg)
	case update.EditedMessage != nil:
		msg := handleEditMessage(update)
		_, err = bot.Send(msg)
	default:
		log.Printf("update type not handled: %+v\n", update)
	}
	return err
}

func handleChatMemberUpdate(update *tgbotapi.Update) tgbotapi.MessageConfig {
	return tgbotapi.NewMessage(barmanID, fmt.Sprintf("Bot status updated: %v", update.MyChatMember))
}

func handleEditMessage(update *tgbotapi.Update) tgbotapi.MessageConfig {
	return tgbotapi.NewMessage(barmanID, fmt.Sprintf("%v edited: [%v]", update.EditedMessage.From.FirstName, update.EditedMessage.Text))
}
