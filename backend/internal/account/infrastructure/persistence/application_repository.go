package persistence

import (
	"context"
	"database/sql"

	"github.com/shige1114/job/backend/internal/account/domain/model"
	"github.com/shige1114/job/backend/internal/account/domain/repository"
	"github.com/shige1114/job/backend/internal/account/infrastructure/persistence/db"
)

type sqliteApplicationRepository struct {
	queries *db.Queries
}

var _ repository.ApplicationRepository = (*sqliteApplicationRepository)(nil)

func NewSqliteApplicationRepository(sqlDB *sql.DB) repository.ApplicationRepository {
	return &sqliteApplicationRepository{
		queries: db.New(sqlDB),
	}
}

func (r *sqliteApplicationRepository) Save(ctx context.Context, app *model.Application) error {
	dto := app.ToDTO()

	_, err := r.queries.GetApplication(ctx, dto.ID)
	if err == sql.ErrNoRows {
		return r.queries.CreateApplication(ctx, db.CreateApplicationParams{
			ID:        dto.ID,
			Email:     dto.Email,
			Code:      dto.Code,
			ExpiresAt: dto.ExpiresAt,
		})
	}
	if err != nil {
		return err
	}

	return r.queries.UpdateApplication(ctx, db.UpdateApplicationParams{
		ID:        dto.ID,
		Email:     dto.Email,
		Code:      dto.Code,
		ExpiresAt: dto.ExpiresAt,
	})
}

func (r *sqliteApplicationRepository) FindByID(ctx context.Context, id model.ID) (*model.Application, error) {
	row, err := r.queries.GetApplication(ctx, id.String())
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return r.mapToDomain(row), nil
}

func (r *sqliteApplicationRepository) FindByEmail(ctx context.Context, email model.Email) (*model.Application, error) {
	row, err := r.queries.GetApplicationByEmail(ctx, email.String())
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return r.mapToDomain(row), nil
}

func (r *sqliteApplicationRepository) mapToDomain(row db.Application) *model.Application {
	id, _ := model.ParseID(row.ID)
	email, _ := model.NewEmail(row.Email)
	return model.NewApplicationFromRepository(
		id,
		email,
		model.NewCode(row.Code),
		row.ExpiresAt,
	)
}
