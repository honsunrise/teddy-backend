package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/zhsyourai/teddy-backend/api/gin-jwt"
)

type Content struct {
	middleware *gin_jwt.JwtMiddleware
}

func NewContentHandler(middleware *gin_jwt.JwtMiddleware) (*Content, error) {
	return &Content{
		middleware: middleware,
	}, nil
}

func (h *Content) Handler(root gin.IRoutes) {
}
