package functions

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/GoogleCloudPlatform/functions-framework-go/functions"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"gitlab.com/gruppi-preparazione/floriande-bot/floriandebot"
	"gitlab.com/gruppi-preparazione/floriande-bot/store"
)

// this init() is called when testing locally with functions-framework-go package
func init() {
	functions.HTTP("BotHandler", BotHandler)
	functions.HTTP("APIHandler", APIHandler)
}

func botSetup() (*tgbotapi.BotAPI, error) {
	bot, err := tgbotapi.NewBotAPI(os.Getenv("TOKEN"))
	if err != nil {
		return nil, err
	}
	bot.Debug = true
	return bot, nil
}

// Entry point for a Telegram message update
func BotHandler(w http.ResponseWriter, r *http.Request) {
	bot, err := botSetup()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Panicln(err)
	}

	floriandebot.Setup(bot, store.SetupFirestore())
	defer floriandebot.CloseDB()

	update, err := bot.HandleUpdate(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		log.Println(err)
		return
	}

	err = floriandebot.HandleUpdate(update)
	if err != nil {
		switch err.(type) {
		case *floriandebot.PlainMsgInGroupError:
			log.Println(err.Error())
		default:
			log.Println(err.Error())
			http.Error(w, err.Error(), http.StatusBadRequest)
		}
	}
}

// Entry point for a call to the API
func APIHandler(w http.ResponseWriter, r *http.Request) {
	bot, err := botSetup()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Panicln(err)
	}

	floriandebot.Setup(bot, store.SetupFirestore())
	defer floriandebot.CloseDB()

	err = floriandebot.HandleAPICall(w, r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		log.Fatal(fmt.Errorf("could not get handle API call: %v", err))
	}
}
