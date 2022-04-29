package controller

import (
	"fmt"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"gitlab.com/gruppi-preparazione/floriande-bot/store"
)

var welcomeText = fmt.Sprintf(`Welcome by Floriande Lounge bar %v\.

Have you already reserved a spot for an upcoming event? You can do so with the `+
	"*`/book`*"+` command\.
Please use the `+"*`/menu`*"+` command to download our latest drink selection\. %v
You can order a drink from here using the `+"`/drink`"+` command, or you can type `+
	`the name of the cocktail if you know it already\. Make sure you spell it correctly\!
Then, you can check if you have any order waiting to be prepared and served`+
	`with the `+"`/orders`"+` command\. %v
Please let [us](tg://user?id=%d) know if you have suggestions for improvement\.
We hope you enjoy you stay\. %v

What would you like to drink today? %v`,
	clinkingGlasses, tumblerGlass, cocktailGlass, barmanID, sun, personTipping)

var helpText = "You can use the command *`/book`* to reserve for an event, " +
	"*`/menu`* to download the digital version of our cocktail menu, " +
	"*`/drink`* to order a cocktail" + ` \(you will be guided through the process\), ` +
	"*`/orders`*" + ` to see the cocktail\(s\) you have ordered and are being mixed\.
Finally, the commands ` + "*`/list`* and *`/serve`*," + ` are reserved for the barman\.`

func handleCommands(update *tgbotapi.Update) error {
	var isBarman bool
	msg := tgbotapi.NewMessage(update.Message.Chat.ID, update.Message.Text)
	switch update.Message.Command() {
	case "start":
		msg.ParseMode = "MarkdownV2"
		msg.Text = welcomeText
	case "help":
		msg.ParseMode = "MarkdownV2"
		msg.Text = helpText
	case "book":
		msg.ReplyMarkup = store.EventsKeyboard
		msg.Text = "We are happy you want to reserve your place with us. Which event would you like to join?"
	case "menu":
		menu1 := tgbotapi.NewInputMediaPhoto(tgbotapi.FileID(store.FreshEasyMenu))
		menu1.Caption = "1. Fresh & Easy menu"
		menu2 := tgbotapi.NewInputMediaPhoto(tgbotapi.FileID(store.ConnoisseurMenu))
		menu2.Caption = "2. Connoisseur menu"
		menuPhotos := []interface{}{menu1, menu2}
		menu := tgbotapi.NewMediaGroup(update.Message.Chat.ID, menuPhotos)
		bot.Send(menu)
		msg.Text = "Here you go, our digital menu. What would you like to drink?"
	case "drink":
		msg.ReplyMarkup = store.CategoriesKeyboard
		msg.Text = fmt.Sprintf("1. let's start with choosing a category %v; what kind of cocktail would you like?", tropicalDrink)
	case "orders":
		orders := store.UserOrders(client, update.Message.From.ID).String()
		msg.Text = fmt.Sprintf("Your order(s):\n%s", orders)
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
		msg.Text = "I don't know that command"
	}
	_, err := bot.Send(msg)
	return err
}

func checkIfBarman(update *tgbotapi.Update, msg tgbotapi.MessageConfig) (bool, tgbotapi.MessageConfig) {
	if update.Message.From.ID != barmanID {
		msg.ParseMode = "MarkdownV2"
		msg.Text = fmt.Sprintf("The `%v` command is reserved for the barman", update.Message.Text)
		return true, msg
	}
	return false, tgbotapi.MessageConfig{}
}
