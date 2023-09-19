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
	APPOINTMENT
	// HAPPENED_IN - MEMORIES
	// RECURRENT
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

// returns true if the input string is a number
func isday(s string) bool {
	_, err := strconv.Atoi(s)
	if err != nil {
		return false
	}
	return true
}

// returns a new calendar with of the specified type
func NewCalendar(ctype int) calendar {
	return calendar{
		Day:          1,
		Month:        time.Now().Month(),
		Year:         time.Now().Year(),
		CalendarType: ctype,
	}
}

// THIS SECTION OF CODE HERE IS SUBJECT TO CHANGE ---------------

func (c *calendar) prevmAppnt() bool {
	if c.Year > time.Now().Year() {
		return true
	}
	return false
}

func (c *calendar) canGetPreviousMonth() bool {
	if c.CalendarType == APPOINTMENT {
		return c.prevmAppnt()
	}
	return true
}

func (c *calendar) nextmBday() bool {
	if c.Month == time.Now().Month() &&
		c.Year == time.Now().Year() {
		return false
	}
	return true
}

func (c *calendar) canGetNextMonth() bool {
	if c.CalendarType == BIRTHDAY {
		return c.nextmBday()
	}
	return true
}

// END OF SECTION --------------------------------------------

// safely gets the next month in the calendar
func (c *calendar) prevm() {
	if int(c.Month) == 1 {
		c.Month = 12
		c.Year--
	} else {
		c.Month--
	}
}

// safely gets the previous month in the calendar
func (c *calendar) nextm() {
	if int(c.Month) == 12 {
		c.Month = 1
		c.Year++
	} else {
		c.Month++
	}
}

// fill the calendar keyboard with the year buttons
func years(c calendar, k keyboard) keyboard {
	yrs := []echotron.InlineKeyboardButton{
		{Text: "<", CallbackData: "prevy"},
		{Text: fmt.Sprint(c.Year), CallbackData: "year"},
		{Text: ">", CallbackData: "nexty"},
	}
	k = append(k, yrs)
	return k
}

// fills the calendar keyboard with the month buttons
func months(c calendar, k keyboard) keyboard {
	mnt := []echotron.InlineKeyboardButton{
		{Text: "<", CallbackData: "prevm"},
		{Text: c.Month.String(), CallbackData: "month"},
		{Text: ">", CallbackData: "nextm"},
	}
	k = append(k, mnt)
	return k
}

// returns a filler button for the calendar inline keyboard
func emptyday(btn button) button {
	btn = append(btn,
		echotron.InlineKeyboardButton{Text: " ", CallbackData: "ignore"})
	return btn
}

// returns a single day button for the calendar inline keyboard
func day(btn button, day int) button {
	btn = append(
		btn,
		echotron.InlineKeyboardButton{
			Text:         fmt.Sprint(day),
			CallbackData: fmt.Sprint(day),
		},
	)
	return btn
}

// fills the calendar inline keyboard with all the days in the month
func days(c calendar, k keyboard) keyboard {
	var (
		tmp []echotron.InlineKeyboardButton
		max = time.Date(c.Year, c.Month+1, 0, 0, 0, 0, 0, time.UTC).Day()
	)

	for days := 1; days <= 31; {
		for row := 0; row < 7; row++ {
			if days > max {
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

// returns a complete calendar inline keyboard that contains year, month and days
func IKbCalendar(c calendar) keyboard {
	var layout keyboard
	layout = years(c, layout)
	layout = months(c, layout)
	layout = days(c, layout)
	return layout
}
