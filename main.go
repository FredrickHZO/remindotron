package main

import (
	_ "embed"
	"fmt"
	"strings"
	"time"

	"github.com/NicoNex/echotron/v3"
)

type bot struct {
	chatID int64
	echotron.API
}

type dateinfo struct {
	date string
	desc string
}

// OLD CODE - TO REMOVE
var infos []dateinfo

var cal calendar

//go:embed token
var token string
var opts echotron.APIResponseMessage

func newBot(chatID int64) echotron.Bot {
	bot := &bot{
		chatID: chatID,
		API:    echotron.NewAPI(token),
	}
	return bot
}

func (b *bot) Update(update *echotron.Update) {
	if update.Message == nil {
		b.handleInlineQueries(update)
		return
	}

	switch {
	case strings.HasPrefix(update.Message.Text, "/start"):
		test := [][]echotron.InlineKeyboardButton{
			{
				{Text: "🎊 Compleanno", CallbackData: "bday"},
				{Text: "📅 Impegno", CallbackData: "date"},
			},
			{
				{Text: "Lista Comp.", CallbackData: "listbday"},
				{Text: "Lista Impe.", CallbackData: "listapp"},
			},
		}
		ok := echotron.InlineKeyboardMarkup{
			InlineKeyboard: test}

		b.SendMessage(
			"Scegli una di queste azioni dalla lista",
			b.chatID,
			&echotron.MessageOptions{ReplyMarkup: ok},
		)
	}
}

func (b *bot) handleCalendar(ctype int) {
	if opts.Result != nil {
		b.DeleteMessage(b.chatID, opts.Result.ID)
	}

	cal = calendar{
		Day:          1,
		Month:        time.Now().Month(),
		Year:         time.Now().Year(),
		CalendarType: ctype,
	}

	rmkup := echotron.InlineKeyboardMarkup{
		InlineKeyboard: generateCalendar(cal),
	}

	str := "Scegli una data dal calendario"
	if ctype == BIRTHDAY {
		str += " per *IL COMPLEANNO* della persona"
	} else {
		str += " per l'*IMPEGNO* che vuoi inserire"
	}

	opts, _ = b.SendMessage(
		str,
		b.chatID,
		&echotron.MessageOptions{
			ReplyMarkup: rmkup,
			ParseMode:   echotron.MarkdownV2,
		},
	)
}

func (b *bot) handleCalendarNextMonth() {
	b.DeleteMessage(b.chatID, opts.Result.ID)

	if int(cal.Month) == 12 {
		cal.Month = 1
	} else {
		cal.Month++
	}
	mkup := echotron.InlineKeyboardMarkup{
		InlineKeyboard: generateCalendar(cal),
	}

	opts, _ = b.SendMessage(
		opts.Result.Text,
		b.chatID,
		&echotron.MessageOptions{
			ReplyMarkup: mkup,
			ParseMode:   echotron.MarkdownV2,
		},
	)
}

func (b *bot) handleInlineQueries(update *echotron.Update) {
	if update.CallbackQuery == nil {
		return
	}

	switch {
	case update.CallbackQuery.Data == "listapp":
		b.SendMessage(
			"*LA LISTA DEGLI IMPEGNI È:* \n"+
				infos[0].date+infos[0].desc,
			b.chatID,
			&echotron.MessageOptions{
				ParseMode: echotron.MarkdownV2,
			},
		)

	case update.CallbackQuery.Data == "date":
		b.handleCalendar(DATE)

	case update.CallbackQuery.Data == "bday":
		b.handleCalendar(BIRTHDAY)

	case update.CallbackQuery.Data == "next":
		b.handleCalendarNextMonth()

	case isNumeric(update.CallbackQuery.Data):
		var (
			year  = fmt.Sprint(cal.Year)
			month = fmt.Sprint(int(cal.Month))
			str   = update.CallbackQuery.Data + "/" + month + "/" + year
		)
		infos = append(infos, dateinfo{date: str, desc: "\tDesc: nessuna"})

		b.SendMessage(
			"Hai selezionato il giorno "+"*"+str+"*"+"\n",
			b.chatID,
			&echotron.MessageOptions{
				ParseMode: echotron.MarkdownV2,
			},
		)
	}

}

func main() {
	dsp := echotron.NewDispatcher(token, newBot)
	dsp.Poll()
}
