package main

import (
	"fmt"
	"strconv"
	"time"

	"github.com/NicoNex/echotron/v3"
)

type keyboard [][]echotron.InlineKeyboardButton
type btnrow []echotron.InlineKeyboardButton

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
	if _, err := strconv.Atoi(s); err != nil {
		return false
	}
	return true
}

// returns a new calendar of the specified type
func NewCalendar(ctype int) calendar {
	return calendar{
		Day:          1,
		Month:        time.Now().Month(),
		Year:         time.Now().Year(),
		CalendarType: ctype,
	}
}

func (c *calendar) canGetPreviousYear() bool {
	switch c.CalendarType {
	case APPOINTMENT:
		if time.Now().Year() == c.Year {
			return false
		}
		return true

	default:
		return true
	}
}

func (c *calendar) canGetNextYear() bool {
	switch c.CalendarType {
	case BIRTHDAY:
		if time.Now().Year() == c.Year {
			return false
		}
		return true

	default:
		return true
	}
}

func (c *calendar) testfunc() bool {
	switch c.CalendarType {
	case APPOINTMENT:
		if c.Year == time.Now().Year()+1 &&
			c.Month < time.Now().Month() {
			return true
		}
		return false

	default:
		if c.Year == time.Now().Year()-1 &&
			c.Month > time.Now().Month() {
			return true
		}
		return false
	}
}

// returns the correct value {true or false} for every type of calendar
// when an action to get the previous month of the one shown in
// the calendar is made
func (c *calendar) canGetPreviousMonth() bool {
	switch c.CalendarType {
	case APPOINTMENT:
		if c.Year > time.Now().Year() {
			return true
		}
		return false

	default:
		return true
	}
}

// returns the correct value {true or false} for every type of calendar
// when an action to get the next month of the one shown in
// the calendar is made
func (c *calendar) canGetNextMonth() bool {
	switch c.CalendarType {
	case BIRTHDAY:
		if c.Month == time.Now().Month() &&
			c.Year == time.Now().Year() {
			return false
		}
		return true

	default:
		return true
	}
}

// safely gets the next month in the calendar, if the month is
// January when this function is called, the year will change accordingly
func (c *calendar) prevm() {
	if c.Month == time.January {
		c.Month = time.December
		c.Year--
	} else {
		c.Month--
	}
}

// safely gets the previous month in the calendar, if the month is
// December when this function is called, the year will change accordingly
func (c *calendar) nextm() {
	if c.Month == time.December {
		c.Month = time.January
		c.Year++
	} else {
		c.Month++
	}
}

// appends the year buttons to the calendar keyboard
func years(c calendar, k keyboard) keyboard {
	yrs := []echotron.InlineKeyboardButton{
		{Text: "<", CallbackData: "prevy"},
		{Text: fmt.Sprint(c.Year), CallbackData: "year"},
		{Text: ">", CallbackData: "nexty"},
	}
	k = append(k, yrs)
	return k
}

// appends the month buttons to the calendar keaybord
func months(c calendar, k keyboard) keyboard {
	mnt := []echotron.InlineKeyboardButton{
		{Text: "<", CallbackData: "prevm"},
		{Text: itmonths[c.Month], CallbackData: "month"},
		{Text: ">", CallbackData: "nextm"},
	}
	k = append(k, mnt)
	return k
}

// appends a single filler day button for a row of the calendar keyboard
func emptyday(btn btnrow) btnrow {
	btn = append(btn,
		echotron.InlineKeyboardButton{Text: " ", CallbackData: "ignore"})
	return btn
}

// appends a single day button for a row of the calendar keyboard
func day(btn btnrow, day int) btnrow {
	btn = append(
		btn,
		echotron.InlineKeyboardButton{
			Text:         fmt.Sprint(day),
			CallbackData: fmt.Sprint(day),
		},
	)
	return btn
}

// appends the correct kind of day button
func putday(max, days int, b btnrow) btnrow {
	if days > max {
		return emptyday(b)
	}
	return day(b, days)
}

// fills the calendar inline keyboard with all the days in the month
func days(c calendar, k keyboard) keyboard {
	max := time.Date(c.Year, c.Month+1, 0, 0, 0, 0, 0, time.UTC).Day()

	for days := 1; days <= 31; {
		tmp := []echotron.InlineKeyboardButton{}
		for row := 0; row < 7; row++ {
			tmp = putday(max, days, tmp)
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
