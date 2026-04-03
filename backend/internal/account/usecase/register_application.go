// Package usecaseはコントローラとの繋ぎ目
package usecase

import (
	"context"
	"crypto/rand"
	"errors"
	"fmt"
	"math/big"
	"time"

	"github.com/google/uuid"
	"github.com/shige1114/job/backend/internal/account/domain/model"
	"github.com/shige1114/job/backend/internal/account/domain/repository"
	"github.com/shige1114/job/backend/internal/account/domain/service"
	"github.com/shige1114/job/backend/internal/account/usecase/command"
)

var (
	ErrEmailAlreadyInUse = errors.New("email already in use")
	ErrInvalidInput      = errors.New("invalid input")
)

type RegisterApplicationUseCase struct {
	repo    repository.ApplicationRepository
	service service.FindSameEmailService
}

func NewRegisterApplicationUseCase(
	repo repository.ApplicationRepository,
	service service.FindSameEmailService,
) *RegisterApplicationUseCase {
	return &RegisterApplicationUseCase{
		repo:    repo,
		service: service,
	}
}

func (u *RegisterApplicationUseCase) Execute(ctx context.Context, input command.RegisterApplicationCommand) (*command.RegisterApplicationResult, error) {
	email, err := model.NewEmail(input.Email)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrInvalidInput, err)
	}

	exists, err := u.service.Exists(ctx, email)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, ErrEmailAlreadyInUse
	}

	uid, err := uuid.NewV7()
	if err != nil {
		return nil, fmt.Errorf("failed to generate uuid: %w", err)
	}

	code, err := GenerateCode()
	if err != nil {
		return nil, fmt.Errorf("failed to generate code: %w", err)
	}

	app := model.NewApplication(
		model.NewID(uid),
		email,
		model.NewCode(code),
		24*time.Hour,
	)

	if err := u.repo.Save(ctx, app); err != nil {
		return nil, fmt.Errorf("failed to save application: %w", err)
	}

	return &command.RegisterApplicationResult{
		ID: app.GetID(),
	}, nil
}

func GenerateCode() (string, error) {
	max := new(big.Int).Exp(big.NewInt(10), big.NewInt(6), nil)
	n, err := rand.Int(rand.Reader, max)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%06d", n), nil
}
