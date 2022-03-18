package main

import (
	"log"
	"os"
	"time"

	tele "gopkg.in/telebot.v3"
)

func main() {
	pref := tele.Settings{
		Token:  os.Getenv("TOKEN"),
		Poller: &tele.LongPoller{Timeout: 10 * time.Second},
	}

	b, err := tele.NewBot(pref)
	if err != nil {
		log.Fatal(err)
		return
	}

	b.Handle("/order", func(c tele.Context) error {
		return c.Send("We will be serving delicous drinks here. Stay tuned!")
	})

	b.Handle("/mojito", func(c tele.Context) error {
		return c.Send("You want to drink a Mojito? Always a good choice!!")
	})

	b.Start()
}
