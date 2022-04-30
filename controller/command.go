package controller

import (
	"fmt"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"gitlab.com/gruppi-preparazione/floriande-bot/store"
)

func handleCommands(update *tgbotapi.Update) error {
	var isBarman bool
	msg := tgbotapi.NewMessage(update.Message.Chat.ID, update.Message.Text)
	switch update.Message.Command() {
	case "start":
		msg.ParseMode = "MarkdownV2"
		msg.Text = fmt.Sprintf(mss[welcomeText][userLanguage(update.Message.From.LanguageCode)],
			clinkingGlasses, tumblerGlass, cocktailGlass, barmanID, sun, personTipping)
	case "help":
		msg.ParseMode = "MarkdownV2"
		msg.Text = mss[helpText][userLanguage(update.Message.From.LanguageCode)]
	case "book":
		msg.ReplyMarkup = store.EventsKeyboard
		msg.Text = mss[chooseEvent][userLanguage(update.Message.From.LanguageCode)]
	case "menu":
		menu1 := tgbotapi.NewInputMediaPhoto(tgbotapi.FileID(store.FreshEasyMenu))
		menu1.Caption = "1. Fresh & Easy menu"
		menu2 := tgbotapi.NewInputMediaPhoto(tgbotapi.FileID(store.ConnoisseurMenu))
		menu2.Caption = "2. Connoisseur menu"
		menuPhotos := []interface{}{menu1, menu2}
		menu := tgbotapi.NewMediaGroup(update.Message.Chat.ID, menuPhotos)
		bot.Send(menu)
		msg.Text = mss[hereTheMenu][userLanguage(update.Message.From.LanguageCode)]
	case "drink":
		msg.ReplyMarkup = store.CategoriesKeyboard
		msg.Text = fmt.Sprintf(mss[chooseCategory][userLanguage(update.Message.From.LanguageCode)], tropicalDrink)
	case "orders":
		orders := store.UserOrders(client, update.Message.From.ID).String()
		msg.Text = fmt.Sprintf(mss[yourOrders][userLanguage(update.Message.From.LanguageCode)], orders)
	case "list":
		if isBarman, msg = checkIfBarman(update, msg); !isBarman {
			break
		}
		orders := store.AllOrders(client).String()
		msg.ParseMode = "MarkdownV2"
		msg.Text = fmt.Sprintf("Outstanding orders:\n"+"```\n"+"%s"+"```", orders)
	case "serve":
		if isBarman, msg = checkIfBarman(update, msg); !isBarman {
			break
		}
		orders := store.AllOrders(client)
		if len(orders) > 0 {
			ordersKeyboard := store.NewServeOrdersKeyboard(orders)
			msg.ReplyMarkup = ordersKeyboard
			msg.Text = fmt.Sprintf("Select order to complete (%2.1d outstanding)", len(orders))
		} else {
			msg.Text = "No waiting orders found"
		}
	default:
		msg.ReplyToMessageID = update.Message.MessageID
		msg.Text = mss[dontKnowCommand][userLanguage(update.Message.From.LanguageCode)]
	}
	_, err := bot.Send(msg)
	return err
}

func checkIfBarman(update *tgbotapi.Update, msg tgbotapi.MessageConfig) (bool, tgbotapi.MessageConfig) {
	if update.Message.From.ID != barmanID {
		msg.ParseMode = "MarkdownV2"
		msg.Text = fmt.Sprintf(mss[reservedForBarman][userLanguage(update.Message.From.LanguageCode)], update.Message.Text)
		return false, msg
	}
	return true, msg
}
