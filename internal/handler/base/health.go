package base

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"teddy-backend/internal/gin_jwt"
)

type Health struct {
	middleware *gin_jwt.JwtMiddleware
}

func NewHealthHandler(middleware *gin_jwt.JwtMiddleware) (*Health, error) {
	return &Health{
		middleware: middleware,
	}, nil
}

func (h *Health) Handler(root gin.IRoutes) {
	root.Any("/", h.ReturnOK)
}

func (h *Health) ReturnOK(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{
		"status": "OK",
	})
}
