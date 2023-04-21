package floriandebot

import (
	"fmt"
	"log"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"gitlab.com/gruppi-preparazione/floriande-bot/store"
)

func handleCallBackQuery(update *tgbotapi.Update) error {
	var err error
	msgText := mss[youChose][userLanguage(update.CallbackQuery.Message.From.LanguageCode)]
	callback := tgbotapi.NewCallback(update.CallbackQuery.ID, fmt.Sprintf(msgText, update.CallbackQuery.Data[1:]))
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
		err = handleOrderCallback(update)
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
		mss[youCanceled][userLanguage(update.CallbackQuery.From.LanguageCode)],
		store.EmptyInlineKeyboard,
	)
	_, err := bot.Request(removeKeyboard)
	return err
}

func handleCategorySelection(update *tgbotapi.Update) error {
	menu := store.CocktailsMenu(client)
	_, cocktailKeyboards := store.NewCocktailKeyboards(menu)
	replaceKeyboard := tgbotapi.NewEditMessageTextAndMarkup(
		update.CallbackQuery.Message.Chat.ID,
		update.CallbackQuery.Message.MessageID,
		mss[chooseDrinkInCategory][userLanguage(update.CallbackQuery.From.LanguageCode)],
		cocktailKeyboards[update.CallbackQuery.Data],
	)
	_, err := bot.Request(replaceKeyboard)
	return err
}

func handleOrderCallback(update *tgbotapi.Update) error {
	drink, category := splitCocktailString(update.CallbackQuery.Data)
	store.AddOrder(client,
		update.CallbackQuery.From.ID,
		update.CallbackQuery.From.FirstName,
		update.CallbackQuery.From.LanguageCode,
		drink,
		category,
	)
	confirm := fmt.Sprintf("%v %s just ordered a %s", bell, update.CallbackQuery.From.FirstName, drink)
	if _, err := bot.Send(tgbotapi.NewMessage(barmanID, confirm)); err != nil {
		log.Fatal(err)
	}
	removeKeyboard := tgbotapi.NewEditMessageTextAndMarkup(
		update.CallbackQuery.Message.Chat.ID,
		update.CallbackQuery.Message.MessageID,
		fmt.Sprintf(mss[orderConfirmation][userLanguage(update.CallbackQuery.From.LanguageCode)], drink, cocktailGlass),
		store.EmptyInlineKeyboard,
	)
	_, err := bot.Request(removeKeyboard)
	return err
}

// Splits cocktail order callbackdata string like "2Cocktail Name1Category Name"
func splitCocktailString(code string) (string, string) {
	chunks := strings.Split(code, "1")
	drink, category := chunks[0][1:], chunks[1]
	return drink, category
}

func handleCloseOrder(update *tgbotapi.Update) error {
	orderID := update.CallbackQuery.Data
	err := closeOrder(orderID)
	if err != nil {
		log.Fatal(fmt.Errorf("could not handle close order: %v", err))
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
