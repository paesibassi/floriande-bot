package controller

import (
	"fmt"
	"log"

	"cloud.google.com/go/firestore"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"gitlab.com/gruppi-preparazione/floriande-bot/store"
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

	barman = 20137373 // admin chat_id
)

var welcomeText = fmt.Sprintf(`Welcome by Floriande Lounge bar %v\.

Please use the `+"*`/menu`*"+` command to download our latest drink selection\. %v
You can order a drink from here using the `+"`/drink`"+` command, and check if you have `+
	`any order waiting to be prepared and served with the `+"`/orders`"+` command\. %v
Please let [us](tg://user?id=20137373) know if you have suggestions for improvement\.
We hope you enjoy you stay\. %v  

What would you like to drink today? %v`,
	clinkingGlasses, tumblerGlass, cocktailGlass, sun, personTipping)

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
		switch update.Message.IsCommand() {
		case true:
			msg = handleCommands(update)
		default:
			msg = handleEcho(update)
		}
		_, err = bot.Send(msg)
	case update.CallbackQuery != nil:
		handleCallBackQuery(update)
	case update.MyChatMember != nil:
		msg := handleChatMemberUpdate(update)
		_, err = bot.Send(msg)
	default:
		err = fmt.Errorf("update type not handled: %+v", update)
	}
	return err
}

func handleCommands(update *tgbotapi.Update) tgbotapi.MessageConfig {
	msg := tgbotapi.NewMessage(update.Message.Chat.ID, update.Message.Text)
	switch update.Message.Command() {
	case "start":
		msg.ParseMode = "MarkdownV2"
		msg.Text = welcomeText
	case "help":
		msg.Text = "I understand `/menu`, `/drink` and `/orders` (and additionally `/list` and `/serve`, if you are a barman)."
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
		msg.Text = fmt.Sprintf("My orders:\n%s", orders)
	case "list":
		orders := store.AllOrders(client).String()
		msg.ParseMode = "MarkdownV2"
		msg.Text = fmt.Sprintf("Outstanding orders:\n"+"```\n"+"%s"+"```", orders)
	case "serve":
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
	return msg
}

func handleCallBackQuery(update *tgbotapi.Update) error {
	var err error
	fmt.Println("Callbackdata: ", update.CallbackQuery.Data)
	callback := tgbotapi.NewCallback(update.CallbackQuery.ID, fmt.Sprintf("You selected %s", update.CallbackQuery.Data[1:]))
	_, err = bot.Request(callback)
	if err != nil {
		return err
	}

	switch update.CallbackQuery.Data[:1] {
	case "0":
		err = handleCancel(update)
	case "1":
		err = handleCategorySelection(update)
	case "2":
		err = handleOrder(update)
	default:
		err = handleCloseOrder(update)
	}

	if err != nil {
		err = fmt.Errorf("callback failed: %s", err)
	}
	return err
}

func handleCancel(update *tgbotapi.Update) error {
	removeKeyboard := tgbotapi.NewEditMessageTextAndMarkup(
		update.CallbackQuery.Message.Chat.ID,
		update.CallbackQuery.Message.MessageID,
		"You canceled your order",
		store.EmptyInlineKeyboard,
	)
	_, err := bot.Request(removeKeyboard)
	return err
}

func handleCategorySelection(update *tgbotapi.Update) error {
	replaceKeyboard := tgbotapi.NewEditMessageTextAndMarkup(
		update.CallbackQuery.Message.Chat.ID,
		update.CallbackQuery.Message.MessageID,
		"2. good choice! Now choose your favorite drink in this category",
		store.CocktailKeyboards[update.CallbackQuery.Data],
	)
	_, err := bot.Request(replaceKeyboard)
	return err
}

func handleOrder(update *tgbotapi.Update) error {
	store.AddOrder(client, update.CallbackQuery.From.ID, update.CallbackQuery.From.FirstName, update.CallbackQuery.Data[1:])
	confirm := fmt.Sprintf("%v %s just ordered a %s", bell, update.CallbackQuery.From.FirstName, update.CallbackQuery.Data[1:])
	if _, err := bot.Send(tgbotapi.NewMessage(barman, confirm)); err != nil {
		log.Fatal(err)
	}
	removeKeyboard := tgbotapi.NewEditMessageTextAndMarkup(
		update.CallbackQuery.Message.Chat.ID,
		update.CallbackQuery.Message.MessageID,
		fmt.Sprintf("A %s is coming soon! %v", update.CallbackQuery.Data[1:], cocktailGlass),
		store.EmptyInlineKeyboard,
	)
	_, err := bot.Request(removeKeyboard)
	return err
}

func handleCloseOrder(update *tgbotapi.Update) error {
	orderID := update.CallbackQuery.Data
	orderPtr, err := store.OrderDetails(client, orderID)
	if err != nil {
		log.Fatal(fmt.Errorf("could not get order detail: %v", err))
	}
	if err := store.CloseOrder(client, orderID); err != nil {
		log.Fatal(fmt.Errorf("could not close order: %v", err))
	}
	confirm := fmt.Sprintf("Your order %s is ready! Enjoy!", orderPtr.CocktailName)
	if _, err := bot.Send(tgbotapi.NewMessage(orderPtr.CustomerID, confirm)); err != nil {
		log.Fatal(err)
	}
	removeKeyboard := tgbotapi.NewEditMessageTextAndMarkup(
		update.CallbackQuery.Message.Chat.ID,
		update.CallbackQuery.Message.MessageID,
		fmt.Sprintf("Closed order: %s", update.CallbackQuery.Data),
		store.EmptyInlineKeyboard,
	)
	_, err = bot.Request(removeKeyboard)
	return err
}

func handleEcho(update *tgbotapi.Update) tgbotapi.MessageConfig {
	return tgbotapi.NewMessage(update.Message.Chat.ID, fmt.Sprintf("%v %s", wineGlass, update.Message.Text))
}

func handleChatMemberUpdate(update *tgbotapi.Update) tgbotapi.MessageConfig {
	return tgbotapi.NewMessage(20137373, fmt.Sprintf("Bot status updated: %v", update.MyChatMember))
}
