package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Router は HTTP ルーターを初期化します
func NewRouter(registerAppHandler *RegisterApplicationHandler) *gin.Engine {
	r := gin.Default()

	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	v1 := r.Group("/api/v1")
	{
		v1.POST("/applications", registerAppHandler.Handle)
	}

	return r
}
