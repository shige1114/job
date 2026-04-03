package repository

import (
	"context"

	"github.com/shige1114/job/backend/internal/account/domain/model"
)

// ApplicationRepository は Application 集約の永続化を担うインターフェースです
type ApplicationRepository interface {
	Save(ctx context.Context, app *model.Application) error
	FindByID(ctx context.Context, id model.ID) (*model.Application, error)
	FindByEmail(ctx context.Context, email model.Email) (*model.Application, error)
}
