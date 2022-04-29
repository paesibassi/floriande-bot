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
		if update.Message.IsCommand() {
			err = handleCommands(update)
		} else {
			err = handlePlainMessage(update)
		}
	case update.CallbackQuery != nil:
		handleCallBackQuery(update)
	case update.MyChatMember != nil:
		err = handleChatMemberUpdate(update)
	case update.EditedMessage != nil:
		err = handleEditMessage(update)
	default:
		log.Printf("update type not handled: %+v\n", update)
	}
	return err
}

func handleChatMemberUpdate(update *tgbotapi.Update) error {
	message := fmt.Sprintf(
		"Bot just joined or left the group '%s', added or removed by [%s](tg://user?id=%d)",
		update.MyChatMember.Chat.Title,
		update.MyChatMember.From.FirstName,
		update.MyChatMember.From.ID,
	)
	msg := tgbotapi.NewMessage(barmanID, message)
	msg.ParseMode = "MarkdownV2"
	_, err := bot.Send(msg)
	return err
}

func handleEditMessage(update *tgbotapi.Update) error {
	msg := tgbotapi.NewMessage(barmanID, fmt.Sprintf("%v edited: [%v]", update.EditedMessage.From.FirstName, update.EditedMessage.Text))
	_, err := bot.Send(msg)
	return err
}
