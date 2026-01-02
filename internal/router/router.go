package router

import (
	"github.com/gin-gonic/gin"

	"short-link-service/internal/handler"
)

func Setup(h *handler.Handler) *gin.Engine {
	r := gin.New()
	r.Use(gin.Recovery())

	r.GET("/healthz", h.Health)
	r.POST("/links", h.CreateLink)
	r.GET("/links/:code/stats", h.Stats)
	r.GET("/:code", h.Redirect)

	return r
}
