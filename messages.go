package main

const (
	BDAY_FUTURE     string = "La persona in questione è nel futuro? scegli una data che sia valida"
	DATE_PAST              = "*Hai la capacità di andare nel passato? inserisci una data valida*"
	CAL_CHOOSE_BDAY        = "Scegli una data dal calendario per *IL COMPLEANNO* della persona"
	CAL_CHOOSE_DATE        = "Scegli una data dal calendario per l'*IMPEGNO* che vuoi inserire"
)

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
