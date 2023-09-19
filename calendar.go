package main

import (
	"fmt"
	"strconv"
	"time"

	"github.com/NicoNex/echotron/v3"
)

type keyboard [][]echotron.InlineKeyboardButton
type button []echotron.InlineKeyboardButton

const (
	BIRTHDAY = iota
	DATE
	// HAPPENED_IN - MEMORIES
)

type calendar struct {
	Day          int
	Month        time.Month
	Year         int
	CalendarType int
}

type list struct {
	events   []calendar
	ListType int
}

func isNumeric(s string) bool {
	_, err := strconv.Atoi(s)
	if err != nil {
		return false
	}
	return true
}

func NewCalendar(ctype int) calendar {
	return calendar{
		Day:          1,
		Month:        time.Now().Month(),
		Year:         time.Now().Year(),
		CalendarType: ctype,
	}
}

// WIP
func (c *calendar) canGetPreviousMonth() bool {
	if c.CalendarType == DATE &&
		c.Month == time.Now().Month() &&
		c.Year == time.Now().Year() {
		return false
	}
	return true
}

// WIP
func (c *calendar) canGetNextMonth() bool {
	if c.CalendarType == BIRTHDAY &&
		c.Month == time.Now().Month() &&
		c.Year == time.Now().Year() {
		return false
	}
	return true
}

func (c *calendar) prevMonth() {
	if int(c.Month) == 1 {
		c.Month = 12
		c.Year--
	} else {
		c.Month -= 1
	}
}

func (c *calendar) nextMonth() {
	if int(c.Month) == 12 {
		c.Month = 1
		c.Year++
	} else {
		c.Month++
	}
}

func genCalendarDate(c calendar) keyboard {
	return [][]echotron.InlineKeyboardButton{
		{
			{Text: ">", CallbackData: "nextyear"},
			{Text: fmt.Sprint(c.Year), CallbackData: "year"},
		},
		{
			{Text: "<", CallbackData: "prev"},
			{Text: c.Month.String(), CallbackData: "month"},
			{Text: ">", CallbackData: "next"},
		},
	}
}

func genCalendarBirthday(c calendar) keyboard {
	return [][]echotron.InlineKeyboardButton{
		{
			{Text: "<", CallbackData: "prevyear"},
			{Text: fmt.Sprint(c.Year), CallbackData: "year"},
		},
		{
			{Text: "<", CallbackData: "prev"},
			{Text: c.Month.String(), CallbackData: "month"},
			{Text: ">", CallbackData: "next"},
		},
	}
}

func appendEmptyDayBtn(btn button) button {
	btn = append(btn,
		echotron.InlineKeyboardButton{Text: " ", CallbackData: "ignore"})
	return btn
}

func appendDayBtn(btn button, j int) button {
	btn = append(
		btn,
		echotron.InlineKeyboardButton{
			Text:         fmt.Sprint(j),
			CallbackData: fmt.Sprint(j),
		},
	)
	return btn
}

func populateDaysBtns(c calendar, k keyboard) keyboard {
	maxdays := time.Date(c.Year, c.Month+1, 0, 0, 0, 0, 0, time.UTC).Day()
	for days := 1; days <= 31; {
		var tmp []echotron.InlineKeyboardButton
		for row := 0; row < 7; row++ {
			if days > maxdays {
				tmp = appendEmptyDayBtn(tmp)
			} else {
				tmp = appendDayBtn(tmp, days)
			}
			days++
		}
		k = append(k, tmp)
	}
	return k
}

func genIKbCalendar(c calendar) keyboard {
	var buttons keyboard
	if c.CalendarType == BIRTHDAY {
		buttons = genCalendarBirthday(c)
	} else {
		buttons = genCalendarDate(c)
	}
	buttons = populateDaysBtns(c, buttons)

	return buttons
}
