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

func isday(s string) bool {
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

// THIS SECTION OF CODE HERE IS SUBJECT TO CHANGE ---------------

func (c *calendar) pmappnttest() bool {
	if c.Year > time.Now().Year() {
		return true
	}
	return false
}

func (c *calendar) canGetPreviousMonth() bool {
	if c.CalendarType == DATE {
		return c.pmappnttest()
	}
	return true
}

func (c *calendar) nmbdtest() bool {
	if c.Month == time.Now().Month() &&
		c.Year == time.Now().Year() {
		return false
	}
	return true
}

func (c *calendar) canGetNextMonth() bool {
	if c.CalendarType == BIRTHDAY {
		return c.nmbdtest()
	}
	return true
}

// END OF SECTION --------------------------------------------

func (c *calendar) prevm() {
	if int(c.Month) == 1 {
		c.Month = 12
		c.Year--
	} else {
		c.Month -= 1
	}
}

func (c *calendar) nextm() {
	if int(c.Month) == 12 {
		c.Month = 1
		c.Year++
	} else {
		c.Month++
	}
}

func getlayout(c calendar) keyboard {
	return [][]echotron.InlineKeyboardButton{
		{
			{Text: "<", CallbackData: "prevy"},
			{Text: fmt.Sprint(c.Year), CallbackData: "year"},
			{Text: ">", CallbackData: "nexty"},
		},
		{
			{Text: "<", CallbackData: "prevm"},
			{Text: c.Month.String(), CallbackData: "month"},
			{Text: ">", CallbackData: "nextm"},
		},
	}
}

func emptyday(btn button) button {
	btn = append(btn,
		echotron.InlineKeyboardButton{Text: " ", CallbackData: "ignore"})
	return btn
}

func day(btn button, j int) button {
	btn = append(
		btn,
		echotron.InlineKeyboardButton{
			Text:         fmt.Sprint(j),
			CallbackData: fmt.Sprint(j),
		},
	)
	return btn
}

func days(k keyboard) keyboard {
	maxdays := time.Date(time.Now().Year(), time.Now().Month()+1, 0, 0, 0, 0, 0, time.UTC).Day()
	for days := 1; days <= 31; {
		var tmp []echotron.InlineKeyboardButton
		for row := 0; row < 7; row++ {
			if days > maxdays {
				tmp = emptyday(tmp)
			} else {
				tmp = day(tmp, days)
			}
			days++
		}
		k = append(k, tmp)
	}
	return k
}

func IKbCalendar(c calendar) keyboard {
	buttons := getlayout(c)
	buttons = days(buttons)
	return buttons
}
