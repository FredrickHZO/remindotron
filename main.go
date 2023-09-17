package main

import (
	_ "embed"
	"log"
	"time"

	"github.com/NicoNex/echotron/v3"
)

type bot struct {
	chatID int64
	echotron.API
}

func (b *bot) Update(update *echotron.Update) {
	qm := "two plus two is four minus one is three quick math"
	if update.Message.Text == "/start" {
		b.SendMessage(qm, b.chatID, nil)
	}
}

//go:embed token
var token string

func newBot(chatID int64) echotron.Bot {
	return &bot{
		chatID,
		echotron.NewAPI(token),
	}
}

func main() {
	dsp := echotron.NewDispatcher(token, newBot)
	for {
		log.Println(dsp.Poll())
		time.Sleep(5 * time.Second)
	}
}
