package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/shige1114/job/backend/internal/account/interface/presenter"
	"github.com/shige1114/job/backend/internal/account/usecase"
	"github.com/shige1114/job/backend/internal/account/usecase/command"
)

type RegisterApplicationHandler struct {
	useCase   *usecase.RegisterApplicationUseCase
	presenter *presenter.RegisterApplicationPresenter
}

func NewRegisterApplicationHandler(
	useCase *usecase.RegisterApplicationUseCase,
	presenter *presenter.RegisterApplicationPresenter,
) *RegisterApplicationHandler {
	return &RegisterApplicationHandler{
		useCase:   useCase,
		presenter: presenter,
	}
}

type registerRequest struct {
	Email string `json:"email" binding:"required,email"`
}

func (h *RegisterApplicationHandler) Handle(c *gin.Context) {
	var req registerRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	cmd := command.RegisterApplicationCommand{
		Email: req.Email,
	}

	result, err := h.useCase.Execute(c.Request.Context(), cmd)
	if err != nil {
		// 出力の関心事はすべて Presenter に任せる
		h.presenter.Error(c, err)
		return
	}

	// 成功時のレスポンスも Presenter に任せる
	h.presenter.Success(c, result)
}
