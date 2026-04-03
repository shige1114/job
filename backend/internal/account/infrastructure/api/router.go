package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/shige1114/job/backend/internal/account/interface/handler"
)

func NewRouter(registerHandler *handler.RegisterApplicationHandler) *gin.Engine {
	r := gin.Default()

	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	v1 := r.Group("/api/v1")
	{
		v1.POST("/applications", registerHandler.Handle)
	}

	return r
}
