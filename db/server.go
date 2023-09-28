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

func decode(val []byte) (cl.Calendar, error) {
	var cal cl.Calendar

	dec := gob.NewDecoder(bytes.NewReader(val))
	if err := dec.Decode(&cal); err != nil {
		return cl.Calendar{}, err
	}

	return cal, nil
}

func PutCalendar(id string, c cl.Calendar) error {
	b, err := encode(c)
	if err != nil {
		return err
	}
	return ccCalendar.Put([]byte(id), b)
}

func GetCalendar(id string) (cl.Calendar, error) {
	b, err := ccCalendar.Get([]byte(id))
	if err != nil {
		return cl.Calendar{}, err
	}
	return decode(b)
}

func DelCalendar(id string) error {
	return ccCalendar.Del([]byte(id))
}
