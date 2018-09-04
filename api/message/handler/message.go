package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/zhsyourai/teddy-backend/api/gin-jwt"
)

type Message struct {
	middleware *gin_jwt.JwtMiddleware
}

func NewMessageHandler(middleware *gin_jwt.JwtMiddleware) (*Message, error) {
	return &Message{
		middleware: middleware,
	}, nil
}

func (h *Message) Handler(root gin.IRoutes) {
	root.GET("")
}

// Message.Register is called by the API as /notify/inbox with post body
func (h *Message) Inbox(ctx *gin.Context) {
}
