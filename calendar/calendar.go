package calendar

import (
	"fmt"
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

type Moment struct {
	Hours   int `json:"hours,omitempty"`
	Minutes int `json:"minutes,omitempty"`
}

type Date struct {
	ID       string     `json:"id"`
	Day      int        `json:"day"`
	Month    time.Month `json:"month"`
	Year     int        `json:"years,omitempty"`
	Hours    Moment     `json:"time,omitempty"`
	DateType int        `json:"datetype"`
}

type Calendar struct {
	Dates    []Date `json:"calendar"`
	ListType int    `json:"calendartype"`
}

// returns a new date of the specified type
func Newdate(ctype int) Date {
	return Date{
		Day:      1,
		Month:    time.Now().Month(),
		Year:     time.Now().Year(),
		DateType: ctype,
	}
}

// checks if the month should be reset to the present one
// if there is a particular (legal) year settting < or > then the
// current one and an illegal month (for the 'actual real' year) is set, when
// the year in the date is changed back to be the same of the
// 'actual real' year, the months should be changed aswell
func (c *Date) ShouldResetMonth() bool {
	switch c.DateType {
	case APPOINTMENT:
		if c.Year == time.Now().Year()+1 &&
			c.Month < time.Now().Month() {
			return true
		}
		return false

	case BIRTHDAY:
		if c.Year == time.Now().Year()-1 &&
			c.Month > time.Now().Month() {
			return true
		}
		return false

	default:
		return false
	}
}

func (c *Date) CanGetPreviousYear() bool {
	switch c.DateType {
	case APPOINTMENT:
		if time.Now().Year() == c.Year {
			return false
		}
		return true

	default:
		return true
	}
}

func (c *Date) CanGetNextYear() bool {
	switch c.DateType {
	case BIRTHDAY:
		if time.Now().Year() == c.Year {
			return false
		}
		return true

	default:
		return true
	}
}

// returns the correct value {true or false} for every type of date
// when an action to get the previous month of the one shown in
// the date is made
func (c *Date) CanGetPreviousMonth() bool {
	switch c.DateType {
	case APPOINTMENT:
		if c.Year > time.Now().Year() {
			return true
		}
		return false

	default:
		return true
	}
}

// returns the correct value {true or false} for every type of date
// when an action to get the next month of the one shown in
// the date is made
func (c *Date) CanGetNextMonth() bool {
	switch c.DateType {
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

// safely gets the next month in the date, if the month is
// January when this function is called, the year will change accordingly
func (c *Date) Prevm() {
	if c.Month == time.January {
		c.Month = time.December
		c.Year--
	} else {
		c.Month--
	}
}

// safely gets the previous month in the date, if the month is
// December when this function is called, the year will change accordingly
func (c *Date) Nextm() {
	if c.Month == time.December {
		c.Month = time.January
		c.Year++
	} else {
		c.Month++
	}
}

// appends the year buttons to the date keyboard
func years(c Date, k keyboard) keyboard {
	yrs := []echotron.InlineKeyboardButton{
		{Text: "<", CallbackData: "prevy"},
		{Text: fmt.Sprint(c.Year), CallbackData: "year"},
		{Text: ">", CallbackData: "nexty"},
	}
	k = append(k, yrs)
	return k
}

// appends the month buttons to the date keaybord
func months(c Date, k keyboard) keyboard {
	mnt := []echotron.InlineKeyboardButton{
		{Text: "<", CallbackData: "prevm"},
		{Text: itmonths[c.Month], CallbackData: "month"},
		{Text: ">", CallbackData: "nextm"},
	}
	k = append(k, mnt)
	return k
}

// appends a single filler day button for a row of the date keyboard
func emptyday(btn btnrow) btnrow {
	btn = append(btn,
		echotron.InlineKeyboardButton{Text: " ", CallbackData: "ignore"})
	return btn
}

// appends a single day button for a row of the date keyboard
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

// fills the date inline keyboard with all the days in the month
func days(c Date, k keyboard) keyboard {
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

// returns a complete date inline keyboard that contains year, month and days
func IKbdate(c Date) keyboard {
	var layout keyboard
	layout = years(c, layout)
	layout = months(c, layout)
	layout = days(c, layout)
	return layout
}
