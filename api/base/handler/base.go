package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/zhsyourai/teddy-backend/api/gin-jwt"
)

type Base struct {
	middleware *gin_jwt.JwtMiddleware
}

func NewBaseHandler(middleware *gin_jwt.JwtMiddleware) (*Base, error) {
	return &Base{
		middleware: middleware,
	}, nil
}

func (h *Base) Handler(root gin.IRoutes) {
}
