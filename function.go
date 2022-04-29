package bot

import (
	"log"
	"net/http"
	"os"

	"github.com/GoogleCloudPlatform/functions-framework-go/functions"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"gitlab.com/gruppi-preparazione/floriande-bot/controller"
	"gitlab.com/gruppi-preparazione/floriande-bot/store"
)

// this init() is called when testing locally with functions-framework-go package
func init() {
	functions.HTTP("Handler", Handler)
}

func Handler(w http.ResponseWriter, r *http.Request) {
	bot, err := tgbotapi.NewBotAPI(os.Getenv("TOKEN"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Panicln(err)
	}
	bot.Debug = true

	controller.Setup(bot, store.SetupFirestore())
	defer controller.CloseDB()

	update, err := bot.HandleUpdate(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		log.Println(err)
		return
	}

	err = controller.HandleUpdate(update)
	if err != nil {
		switch err.(type) {
		case *controller.PlainMsgInGroupError:
			log.Println(err.Error())
		default:
			http.Error(w, err.Error(), http.StatusBadRequest)
			log.Println(err)
		}
	}
}
