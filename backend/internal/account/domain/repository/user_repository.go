package repository

import (
	"context"


	"github.com/shige1114/job/backend/internal/account/domain/model"
)

// UserRepository は User 集約の永続化を担うインターフェースです
type UserRepository interface {
	FindByEmail(ctx context.Context, email model.Email) (bool, error)
}
