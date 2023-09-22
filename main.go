package main

import (
	_ "embed"
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
	dt    date
	msg   echotron.APIResponseMessage
)

// returns a new bot
func newBot(chatID int64) echotron.Bot {
	return &bot{
		chatID: chatID,
		API:    echotron.NewAPI(token),
	}
}

// returns the date inline keayboard markup
func cikm() echotron.InlineKeyboardMarkup {
	return echotron.InlineKeyboardMarkup{
		InlineKeyboard: IKbdate(dt),
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
	dt = Newdate(ctype)
	b.send(introMsg(ctype))
}

// sends a new date inline keyboard with the next month
// if it's possible, otherwise with the month that was
// previously set and shows the correct warning message
func (b *bot) handledateNextMonth() {
	if !dt.canGetNextMonth() {
		b.send(errMsg(dt.dateType))
	} else {
		dt.nextm()
		b.editmsg(introMsg(dt.dateType))
	}
}

// sends a new date inline keyboard with the previous month
// if it's possible, otherwise with the month that was
// previously set and shows the correct warning message
func (b *bot) handledatePrevMonth() {
	if !dt.canGetPreviousMonth() {
		b.send(errMsg(dt.dateType))
	} else {
		dt.prevm()
		b.editmsg(introMsg(dt.dateType))
	}
}

// WIP
func (b *bot) handleNextYear() {
	if !dt.canGetNextYear() {
		b.send(errMsg(dt.dateType))
	} else {
		if dt.shouldResetMonth() {
			dt.Month = time.Now().Month()
		}
		dt.Year++
		b.editmsg(introMsg(dt.dateType))
	}
}

func (b *bot) handlePrevYear() {
	if !dt.canGetPreviousYear() {
		b.send(errMsg(dt.dateType))
	} else {
		if dt.shouldResetMonth() {
			dt.Month = time.Now().Month()
		}
		dt.Year--
		b.editmsg(introMsg(dt.dateType))
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
		b.handledate(APPOINTMENT)

	case update.CallbackQuery.Data == "bday":
		b.handledate(BIRTHDAY)

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
