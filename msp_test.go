package msp

import (
	"errors"
)

type Database map[string][]byte

func (d Database) Get(name string) ([]byte, error) {
	out, ok := d[name]

	if ok {
		return out, nil
	} else {
		return []byte(""), errors.New("Not found!")
	}
}
