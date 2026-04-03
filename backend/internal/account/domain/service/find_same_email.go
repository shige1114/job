package service

import (
	"context"

	"github.com/shige1114/job/backend/internal/account/domain/model"
)

// FindSameEmailService は、メールアドレスの重複を確認するためのインターフェースです
type FindSameEmailService interface {
	Exists(ctx context.Context, email model.Email) (bool, error)
}
