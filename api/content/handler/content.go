package handler

import (
	"github.com/gin-gonic/gin"
)

type Content struct {
}

func NewContentHandler() (*Content, error) {
	return &Content{}, nil
}

func (h *Content) Handler(root gin.IRoutes) {
}
