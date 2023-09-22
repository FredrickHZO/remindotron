package main

import (
	_ "embed"
	cl "remindotron/calendar"
	"strconv"
	"strings"
	"time"

	"github.com/NicoNex/echotron/v3"
)

type bot struct {
	chatID int64
	echotron.API
}

var (
	//go:embed token
	token string
	date  cl.Date
	msg   echotron.APIResponseMessage
)

// returns a new bot
func newBot(chatID int64) echotron.Bot {
	return &bot{
		chatID: chatID,
		API:    echotron.NewAPI(token),
	}
}

// returns true if the input string is a number
func isday(s string) bool {
	if _, err := strconv.Atoi(s); err != nil {
		return false
	}
	return true
}

// returns the date inline keayboard markup
func cikm() echotron.InlineKeyboardMarkup {
	return echotron.InlineKeyboardMarkup{
		InlineKeyboard: cl.IKbdate(date),
	}
}

// edits the message with a new date inline keyboard
func (b *bot) editmsg(str string) {
	msg, _ = b.EditMessageText(
		str,
		echotron.NewMessageID(b.chatID, msg.Result.ID),
		&echotron.MessageTextOptions{
			ParseMode:   echotron.MarkdownV2,
			ReplyMarkup: cikm(),
		},
	)
}

// sends a message containing the input string and the date
// inline keyboard markup.
// if in the bots previous messages there is a date, it
// gets deleted.
func (b *bot) send(str string) {
	if msg.Result != nil {
		b.DeleteMessage(b.chatID, msg.Result.ID)
	}
	msg, _ = b.SendMessage(
		str,
		b.chatID,
		&echotron.MessageOptions{
			ParseMode:   echotron.MarkdownV2,
			ReplyMarkup: cikm(),
		},
	)
}

// handles the creation of a new date of the specified type
func (b *bot) handledate(ctype int) {
	date = cl.Newdate(ctype)
	b.send(cl.IntroMsg(ctype))
}

// sends a new date inline keyboard with the next month
// if it's possible, otherwise with the month that was
// previously set and shows the correct warning message
func (b *bot) handledateNextMonth() {
	if !date.CanGetNextMonth() {
		b.send(cl.ErrMsg(date.DateType))
	} else {
		date.Nextm()
		b.editmsg(cl.IntroMsg(date.DateType))
	}
}

// sends a new date inline keyboard with the previous month
// if it's possible, otherwise with the month that was
// previously set and shows the correct warning message
func (b *bot) handledatePrevMonth() {
	if !date.CanGetPreviousMonth() {
		b.send(cl.ErrMsg(date.DateType))
	} else {
		date.Prevm()
		b.editmsg(cl.IntroMsg(date.DateType))
	}
}

// WIP
func (b *bot) handleNextYear() {
	if !date.CanGetNextYear() {
		b.send(cl.ErrMsg(date.DateType))
	} else {
		if date.ShouldResetMonth() {
			date.Month = time.Now().Month()
		}
		date.Year++
		b.editmsg(cl.IntroMsg(date.DateType))
	}
}

func (b *bot) handlePrevYear() {
	if !date.CanGetPreviousYear() {
		b.send(cl.ErrMsg(date.DateType))
	} else {
		if date.ShouldResetMonth() {
			date.Month = time.Now().Month()
		}
		date.Year--
		b.editmsg(cl.IntroMsg(date.DateType))
	}
}

// handles every interaction with the inline keyboard the bot shows
func (b *bot) handleInlineQueries(update *echotron.Update) {
	switch {
	// no db coded yet, can't generate lists of dates
	case update.CallbackQuery.Data == "listappnt":
		b.SendMessage("Funzione non ancora implementata", b.chatID, nil)

	case update.CallbackQuery.Data == "listbday":
		b.SendMessage("Funzione non ancora implementata", b.chatID, nil)

	case update.CallbackQuery.Data == "appnt":
		b.handledate(cl.APPOINTMENT)

	case update.CallbackQuery.Data == "bday":
		b.handledate(cl.BIRTHDAY)

	case update.CallbackQuery.Data == "nextm":
		b.handledateNextMonth()

	case update.CallbackQuery.Data == "prevm":
		b.handledatePrevMonth()

	case update.CallbackQuery.Data == "nexty":
		b.handleNextYear()

	case update.CallbackQuery.Data == "prevy":
		b.handlePrevYear()

	case isday(update.CallbackQuery.Data):
		// WIP, no db to store the chosen date
	}
}

// handles user input
func (b *bot) Update(update *echotron.Update) {
	if update.Message == nil && update.CallbackQuery != nil {
		b.handleInlineQueries(update)
		return
	}
	switch {
	case strings.HasPrefix(update.Message.Text, "/start"):
		ikb := [][]echotron.InlineKeyboardButton{
			{
				{Text: "🎊 Compleanno", CallbackData: "bday"},
				{Text: "📅 Impegno", CallbackData: "appnt"},
			},
			{
				{Text: "Lista Comp.", CallbackData: "listbday"},
				{Text: "Lista Impe.", CallbackData: "listappnt"},
			},
		}
		ikm := echotron.InlineKeyboardMarkup{
			InlineKeyboard: ikb}
		b.SendMessage(
			"Seleziona una di queste azioni disponibili",
			b.chatID,
			&echotron.MessageOptions{ReplyMarkup: ikm},
		)
	}
}

func main() {
	dsp := echotron.NewDispatcher(token, newBot)
	dsp.Poll()
}
