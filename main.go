package main

import (
	_ "embed"
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

// returns a new bot
func newBot(chatID int64) echotron.Bot {
	return &bot{
		chatID: chatID,
		API:    echotron.NewAPI(token),
	}
}

// returns the calendar inline keayboard markup
func cikm() echotron.InlineKeyboardMarkup {
	return echotron.InlineKeyboardMarkup{
		InlineKeyboard: IKbCalendar(cal),
	}
}

// sends a message containing the input string and the calendar
// inline keyboard markup.
// if in the bots previous messages there is a calendar, it
// gets deleted.
func (b *bot) editmsg(str string) {
	b.EditMessageText(
		str,
		echotron.NewMessageID(b.chatID, opts.Result.ID),
		&echotron.MessageTextOptions{
			ParseMode:   echotron.MarkdownV2,
			ReplyMarkup: cikm(),
		},
	)
}

func (b *bot) editmarkup() {
	b.EditMessageReplyMarkup(
		echotron.NewMessageID(b.chatID, opts.Result.ID),
		&echotron.MessageReplyMarkup{
			ReplyMarkup: cikm(),
		},
	)
}

// handles the creation of a new calendar of the specified type
func (b *bot) handleCalendar(ctype int) {
	if opts.Result != nil {
		b.DeleteMessage(b.chatID, opts.Result.ID)
	}
	cal = NewCalendar(ctype)
	opts, _ = b.SendMessage(
		introMsg(cal.CalendarType),
		b.chatID,
		&echotron.MessageOptions{
			ParseMode:   echotron.MarkdownV2,
			ReplyMarkup: cikm(),
		},
	)
}

// sends a new calendar inline keyboard with the next month
// if it's possible, otherwise with the month that was
// previously set and shows the correct warning message
func (b *bot) handleCalendarNextMonth() {
	if !cal.canGetNextMonth() {
		b.editmsg(errMsg(cal.CalendarType))
	} else {
		cal.nextm()
		b.editmarkup()
	}
}

// sends a new calendar inline keyboard with the previous month
// if it's possible, otherwise with the month that was
// previously set and shows the correct warning message
func (b *bot) handleCalendarPrevMonth() {
	if !cal.canGetPreviousMonth() {
		b.editmsg(errMsg(cal.CalendarType))
	} else {
		cal.prevm()
		b.editmarkup()
	}
}

// WIP
func (b *bot) handleNextYear() {
	b.editmsg(introMsg(cal.CalendarType))
}

// handles every interaction with the inline keyboard the bot shows
func (b *bot) handleInlineQueries(update *echotron.Update) {
	switch {
	case update.CallbackQuery.Data == "listappnt":
		b.SendMessage("Funzione non ancora implementata", b.chatID, nil)

	case update.CallbackQuery.Data == "listbday":
		b.SendMessage("Funzione non ancora implementata", b.chatID, nil)

	case update.CallbackQuery.Data == "appnt":
		b.handleCalendar(APPOINTMENT)

	case update.CallbackQuery.Data == "bday":
		b.handleCalendar(BIRTHDAY)

	case update.CallbackQuery.Data == "nextm":
		b.handleCalendarNextMonth()

	case update.CallbackQuery.Data == "prevm":
		b.handleCalendarPrevMonth()

	case update.CallbackQuery.Data == "nexty":
		b.handleNextYear()

	case update.CallbackQuery.Data == "prevy":
		b.SendMessage("Funzione non ancora implementata", b.chatID, nil)

	case isday(update.CallbackQuery.Data):
		// WIP
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
