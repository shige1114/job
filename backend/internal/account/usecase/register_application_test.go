package usecase

import (
	"context"
	"errors"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/shige1114/job/backend/internal/account/usecase/command"
	"github.com/shige1114/job/backend/internal/account/usecase/mock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRegisterApplicationUseCase_Execute(t *testing.T) {
	ctx := context.Background()

	t.Run("成功: 新しい申請を登録できる", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		repo := mock.NewMockApplicationRepository(ctrl)
		svc := mock.NewMockFindSameEmailService(ctrl)

		svc.EXPECT().
			Exists(ctx, gomock.Any()).
			Return(false, nil)

		repo.EXPECT().
			Save(ctx, gomock.Any()).
			Return(nil)

		uc := NewRegisterApplicationUseCase(repo, svc)
		input := command.RegisterApplicationCommand{
			Email: "test@example.com",
		}

		result, err := uc.Execute(ctx, input)

		require.NoError(t, err)
		require.NotNil(t, result)

		_, uuidErr := uuid.Parse(result.ID)
		assert.NoError(t, uuidErr)
	})

	t.Run("失敗: 無効なメールアドレス", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		repo := mock.NewMockApplicationRepository(ctrl)
		svc := mock.NewMockFindSameEmailService(ctrl)

		uc := NewRegisterApplicationUseCase(repo, svc)
		input := command.RegisterApplicationCommand{
			Email: "invalid-email",
		}

		_, err := uc.Execute(ctx, input)
		assert.ErrorIs(t, err, ErrInvalidInput)
	})

	t.Run("失敗: メールアドレスが既に使用されている", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		repo := mock.NewMockApplicationRepository(ctrl)
		svc := mock.NewMockFindSameEmailService(ctrl)

		svc.EXPECT().
			Exists(ctx, gomock.Any()).
			Return(true, nil)

		uc := NewRegisterApplicationUseCase(repo, svc)
		input := command.RegisterApplicationCommand{
			Email: "duplicate@example.com",
		}

		_, err := uc.Execute(ctx, input)
		assert.ErrorIs(t, err, ErrEmailAlreadyInUse)
	})

	t.Run("失敗: リポジトリの保存エラー", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		saveErr := errors.New("db error")
		repo := mock.NewMockApplicationRepository(ctrl)
		svc := mock.NewMockFindSameEmailService(ctrl)

		svc.EXPECT().
			Exists(ctx, gomock.Any()).
			Return(false, nil)

		repo.EXPECT().
			Save(ctx, gomock.Any()).
			Return(saveErr)

		uc := NewRegisterApplicationUseCase(repo, svc)
		input := command.RegisterApplicationCommand{
			Email: "test@example.com",
		}

		_, err := uc.Execute(ctx, input)
		assert.ErrorIs(t, err, saveErr)
	})
}
