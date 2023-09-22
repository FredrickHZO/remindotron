package main

import "time"

const (
	BDAY_FUTURE      string = "*Devo pensare che sei nel futuro? scegli una data che sia valida*"
	APPNT_PAST              = "*Hai la capacità di viaggiare nel tempo? inserisci una data valida*"
	CAL_CHOOSE_BDAY         = "Scegli una data dal calendario per *IL COMPLEANNO* della persona"
	CAL_CHOOSE_APPNT        = "Scegli una data dal calendario per l'*IMPEGNO* che vuoi inserire"
)

var itmonths = map[time.Month]string{
	time.January:   "Gennaio",
	time.February:  "Febbraio",
	time.March:     "Marzo",
	time.April:     "Aprile",
	time.May:       "Maggio",
	time.June:      "Giugno",
	time.July:      "Luglio",
	time.August:    "Agosto",
	time.September: "Settembre",
	time.October:   "Ottobre",
	time.November:  "Novembre",
	time.December:  "Dicembre",
}

// returns the correct intro message the bot shows
// according to the type of the calendar when you
// select one of the calendar options
func introMsg(ctype int) string {
	switch ctype {
	case BIRTHDAY:
		return CAL_CHOOSE_BDAY
	case APPOINTMENT:
		return CAL_CHOOSE_APPNT
	default:
		return " "
	}
}

// returns the correct error message the bot shows
// according to the type of the calendar when the user
// does an illegal action
func errMsg(ctype int) string {
	switch ctype {
	case BIRTHDAY:
		return BDAY_FUTURE
	case APPOINTMENT:
		return APPNT_PAST
	default:
		return " "
	}
}
