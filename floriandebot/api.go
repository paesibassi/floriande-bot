package floriandebot

import (
	"fmt"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"gitlab.com/gruppi-preparazione/floriande-bot/store"
)

// Accepts an orderID, looks up the order detail in the DB and closes the order
func closeOrder(orderID string) error {
	orderPtr, err := store.OrderDetails(client, orderID)
	if err != nil {
		return fmt.Errorf("could not get order detail: %v", err)
	}
	if err := store.CloseOrder(client, orderID); err != nil {
		return fmt.Errorf("could not close order: %v", err)
	}
	msgText := mss[orderReady][userLanguage(orderPtr.CustomerLang)]
	confirm := fmt.Sprintf(msgText, orderPtr.CocktailName)
	_, err = bot.Send(tgbotapi.NewMessage(orderPtr.CustomerID, confirm))
	return err
}
