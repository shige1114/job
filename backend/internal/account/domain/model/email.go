package model

import (
	"errors"
	"net/mail"
)

var ErrInvalidEmail = errors.New("invalid email format")

// Email はバリデーション済みのメールアドレスです
type Email struct {
	value string
}

func NewEmail(v string) (Email, error) {
	addr, err := mail.ParseAddress(v)
	if err != nil {
		return Email{}, ErrInvalidEmail
	}
	if addr.Address != v {
		return Email{}, ErrInvalidEmail
	}
	return Email{value: addr.Address}, nil
}

func (e Email) String() string {
	return e.value
}
