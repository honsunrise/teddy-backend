package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/micro/go-log"
	"github.com/zhsyourai/teddy-backend/api/gin-jwt"
)

type Message struct {
	middleware *gin_jwt.JwtMiddleware
}

func NewMessageHandler(middleware *gin_jwt.JwtMiddleware) *Message {
	return &Message{}
}

func (h *Message) Handler(root gin.IRoutes) {
}

// Message.Register is called by the API as /notify/inbox with post body
func (h *Message) Inbox(ctx *gin.Context) error {
	log.Log("Received Message.Register request")
	return nil
}
