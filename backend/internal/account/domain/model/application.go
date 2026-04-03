// Package model はアカウント作成の申請（Application）集約を定義します
package model

import (
	"errors"
	"time"
)

var (
	ErrExpired     = errors.New("application has expired")
	ErrInvalidCode = errors.New("invalid confirmation code")
)

// Application はアカウント登録の申請を表す集約ルートです
type Application struct {
	id        ID
	code      Code
	email     Email
	password  string
	expiresAt time.Time
}

// NewApplication は新しい申請を登録します (コマンド: 申請の登録)
func NewApplication(id ID, email Email, code Code, ttl time.Duration) *Application {
	return &Application{
		id:        id,
		email:     email,
		code:      code,
		expiresAt: time.Now().Add(ttl),
	}
}

// NewApplicationFromRepository はリポジトリからの再構成用です
func NewApplicationFromRepository(id ID, email Email, code Code, expiresAt time.Time) *Application {
	return &Application{
		id:        id,
		email:     email,
		code:      code,
		expiresAt: expiresAt,
	}
}

// ConfirmEmail は提供されたコードを検証し、メールアドレスの所有を確認します
func (a *Application) ConfirmEmail(code Code) error {
	if a.IsExpired() {
		return ErrExpired
	}
	if a.code.String() != code.String() {
		return ErrInvalidCode
	}
	return nil
}

// RegisterPassword は申請に対してパスワードを設定します
func (a *Application) RegisterPassword(password string) error {
	if a.IsExpired() {
		return ErrExpired
	}
	a.password = password
	return nil
}

// IsExpired は申請が有効期限切れかどうかを判定します
func (a *Application) IsExpired() bool {
	return time.Now().After(a.expiresAt)
}

// Getters

func (a *Application) GetID() string {
	return a.id.String()
}

func (a *Application) GetEmail() string {
	return a.email.String()
}

func (a *Application) ToDTO() ApplicationDTO {
	return ApplicationDTO{
		ID:        a.id.String(),
		Email:     a.email.String(),
		Code:      a.code.String(),
		ExpiresAt: a.expiresAt,
	}
}
