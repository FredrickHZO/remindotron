package main

import (
	"bytes"
	"encoding/gob"
)

var ccCalendar Cache

func encode(cl calendar) ([]byte, error) {
	var buf bytes.Buffer

	enc := gob.NewEncoder(&buf)
	if err := enc.Encode(cl); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}
