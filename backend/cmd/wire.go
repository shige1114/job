//go:build wireinject

package main

import (
	"database/sql"

	"github.com/gin-gonic/gin"
	"github.com/google/wire"
	"github.com/shige1114/job/backend/internal/account/infrastructure/api"
	"github.com/shige1114/job/backend/internal/account/infrastructure/persistence"
	"github.com/shige1114/job/backend/internal/account/interface/handler"
	"github.com/shige1114/job/backend/internal/account/interface/presenter"
	"github.com/shige1114/job/backend/internal/account/usecase"
)

// InitializeRouter は依存関係を解決してルーターを初期化します
func InitializeRouter(db *sql.DB) *gin.Engine {
	wire.Build(
		// リポジトリ
		persistence.NewSqliteApplicationRepository,
		persistence.NewSqliteUserRepository,

		// ドメインサービス
		persistence.NewFindSameEmailService,

		// ユースケース
		usecase.NewRegisterApplicationUseCase,

		// プレゼンター・ハンドラー
		presenter.NewRegisterApplicationPresenter,
		handler.NewRegisterApplicationHandler,

		// API (Router)
		api.NewRouter,
	)
	return nil
}
