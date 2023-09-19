package main

const (
	BDAY_FUTURE     string = "*Devo pensare che sei nel futuro? scegli una data che sia valida*"
	DATE_PAST              = "*Hai la capacità di viaggiare nel tempo? inserisci una data valida*"
	CAL_CHOOSE_BDAY        = "Scegli una data dal calendario per *IL COMPLEANNO* della persona"
	CAL_CHOOSE_DATE        = "Scegli una data dal calendario per l'*IMPEGNO* che vuoi inserire"
)

// returns the correct intro message the bot shows
// according to the type of the calendar when you
// select one of the calendar options
func introMsg(ctype int) string {
	switch ctype {
	case BIRTHDAY:
		return CAL_CHOOSE_BDAY
	case DATE:
		return CAL_CHOOSE_DATE
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
	case DATE:
		return DATE_PAST
	default:
		return " "
	}
}
