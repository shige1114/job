package persistence

import (
	"context"
	"database/sql"

	"github.com/shige1114/job/backend/internal/account/domain/model"
	"github.com/shige1114/job/backend/internal/account/domain/repository"
	"github.com/shige1114/job/backend/internal/account/infrastructure/persistence/db"
)

type sqliteUserRepository struct {
	queries *db.Queries
}

var _ repository.UserRepository = (*sqliteUserRepository)(nil)

func NewSqliteUserRepository(sqlDB *sql.DB) repository.UserRepository {
	return &sqliteUserRepository{
		queries: db.New(sqlDB),
	}
}

func (r *sqliteUserRepository) FindByEmail(ctx context.Context, email model.Email) (bool, error) {
	_, err := r.queries.GetUserByEmail(ctx, email.String())
	if err != nil {
		if err == sql.ErrNoRows {
			return false, nil
		}
		return false, err
	}
	return true, nil
}
