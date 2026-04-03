package presenter

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/shige1114/job/backend/internal/account/usecase"
	"github.com/shige1114/job/backend/internal/account/usecase/command"
)

// RegisterApplicationPresenter は申請登録のレスポンスを構築します
type RegisterApplicationPresenter struct{}

func NewRegisterApplicationPresenter() *RegisterApplicationPresenter {
	return &RegisterApplicationPresenter{}
}

// Success は成功時のレスポンスを返します (201 Created)
func (p *RegisterApplicationPresenter) Success(c *gin.Context, result *command.RegisterApplicationResult) {
	c.JSON(http.StatusCreated, result)
}

// Error はエラー時のレスポンスを返します (400, 409, 500)
func (p *RegisterApplicationPresenter) Error(c *gin.Context, err error) {
	switch {
	case errors.Is(err, usecase.ErrInvalidInput):
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	case errors.Is(err, usecase.ErrEmailAlreadyInUse):
		c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
	default:
		// 内部エラーはセキュリティのため詳細を伏せる
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
	}
}
