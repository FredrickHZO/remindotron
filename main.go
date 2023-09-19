package main

import (
	_ "embed"
	"fmt"
	"strconv"
	"strings"

	"github.com/NicoNex/echotron/v3"
)

type bot struct {
	chatID int64
	echotron.API
}

var (
	//go:embed token
	token string
	cal   calendar
	opts  echotron.APIResponseMessage
)

func newBot(chatID int64) echotron.Bot {
	return &bot{
		chatID: chatID,
		API:    echotron.NewAPI(token),
	}
}

func ikm() echotron.InlineKeyboardMarkup {
	return echotron.InlineKeyboardMarkup{
		InlineKeyboard: generateCalendar(cal),
	}
}

func (b *bot) sendkb(str string) {
	if opts.Result != nil {
		b.DeleteMessage(b.chatID, opts.Result.ID)
	}
	opts, _ = b.SendMessage(
		str,
		b.chatID,
		&echotron.MessageOptions{
			ReplyMarkup: ikm(),
			ParseMode:   echotron.MarkdownV2,
		},
	)
}

func (b *bot) handleCalendar(ctype int) {
	cal = NewCalendar(ctype)
	b.sendkb(introMsg(cal.CalendarType))
}

func (b *bot) handleCalendarNextMonth() {
	if !cal.canGetNextMonth() {
		b.sendkb(errMsg(cal.CalendarType))
	} else {
		cal.nextMonth()
		b.sendkb(introMsg(cal.CalendarType))
	}
}

func (b *bot) handleCalendarPrevMonth() {
	if !cal.canGetPreviousMonth() {
		b.sendkb(errMsg(cal.CalendarType))
	} else {
		cal.prevMonth()
		b.sendkb(introMsg(cal.CalendarType))
	}
}

func (b *bot) handleNextYear() {
	// WIP
	b.sendkb(introMsg(cal.CalendarType))
}

func (b *bot) handleInlineQueries(update *echotron.Update) {
	switch {
	case update.CallbackQuery.Data == "listapp":
		b.SendMessage("Funzione non ancora implementata", b.chatID, nil)

	case update.CallbackQuery.Data == "date":
		b.handleCalendar(DATE)

	case update.CallbackQuery.Data == "bday":
		b.handleCalendar(BIRTHDAY)

	case update.CallbackQuery.Data == "next":
		b.handleCalendarNextMonth()

	case update.CallbackQuery.Data == "prev":
		b.handleCalendarPrevMonth()

	case update.CallbackQuery.Data == "nextyear":
		b.handleNextYear()

	case isNumeric(update.CallbackQuery.Data):
		cal.Day, _ = strconv.Atoi(update.CallbackQuery.Data)
		str := fmt.Sprint(cal.Day) + "/" + cal.Month.String() + "/" + fmt.Sprint(cal.Year)

		b.SendMessage(
			"Hai selezionato il giorno "+"*"+str+"*"+"\n",
			b.chatID,
			&echotron.MessageOptions{
				ParseMode: echotron.MarkdownV2,
			},
		)
	}
}

func (b *bot) Update(update *echotron.Update) {
	if update.Message == nil && update.CallbackQuery != nil {
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

func main() {
	dsp := echotron.NewDispatcher(token, newBot)
	dsp.Poll()
}
