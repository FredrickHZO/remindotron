package db

import (
	"bytes"
	"encoding/gob"
	cl "remindotron/calendar"
)

var ccCalendar Cache

func encode(c cl.Calendar) ([]byte, error) {
	var buf bytes.Buffer

	enc := gob.NewEncoder(&buf)
	if err := enc.Encode(c); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}
