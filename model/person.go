package model

import (
	"errors"
	"strings"
)

type Person struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	Age  int    `json:"age"`
}

func (p Person) Valid() error {

	msg := ""
	if len(strings.TrimSpace(p.ID)) == 0 {
		msg += "id cannot be empty, "
	}

	if len(strings.TrimSpace(p.Name)) == 0 {
		msg += "name cannot be empty, "
	}
	strings.TrimSuffix(msg, ", ")

	if len(msg) > 0 {
		return errors.New(msg)
	}
	return nil
}
