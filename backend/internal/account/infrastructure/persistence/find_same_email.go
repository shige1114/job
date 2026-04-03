package persistence

import (
	"context"
	"fmt"

	"github.com/shige1114/job/backend/internal/account/domain/model"
	"github.com/shige1114/job/backend/internal/account/domain/repository"
	"github.com/shige1114/job/backend/internal/account/domain/service"
)

type sqliteFindSameEmailService struct {
	appRepo  repository.ApplicationRepository
	userRepo repository.UserRepository
}

// インターフェースを満たしているか確認
var _ service.FindSameEmailService = (*sqliteFindSameEmailService)(nil)

func NewFindSameEmailService(
	appRepo repository.ApplicationRepository,
	userRepo repository.UserRepository,
) service.FindSameEmailService {
	return &sqliteFindSameEmailService{
		appRepo:  appRepo,
		userRepo: userRepo,
	}
}

func (s *sqliteFindSameEmailService) Exists(ctx context.Context, email model.Email) (bool, error) {
	app, err := s.appRepo.FindByEmail(ctx, email)
	if err != nil {
		return false, fmt.Errorf("failed to search in application repository: %w", err)
	}
	if app != nil {
		return true, nil
	}

	existsInUser, err := s.userRepo.FindByEmail(ctx, email)
	if err != nil {
		return false, fmt.Errorf("failed to search in user repository: %w", err)
	}
	if existsInUser {
		return true, nil
	}

	return false, nil
}
