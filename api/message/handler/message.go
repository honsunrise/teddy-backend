package handler

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/prometheus/common/log"
	"github.com/zhsyourai/teddy-backend/api/gin-jwt"
	"github.com/zhsyourai/teddy-backend/api/message/client"
	"github.com/zhsyourai/teddy-backend/common/errors"
	"github.com/zhsyourai/teddy-backend/message/proto"
	"net/http"
)

type Message struct {
	middleware *gin_jwt.JwtMiddleware
}

func NewMessageHandler(middleware *gin_jwt.JwtMiddleware) (*Message, error) {
	return &Message{
		middleware: middleware,
	}, nil
}

var wsUpgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func (h *Message) Handler(root gin.IRoutes) {
	root.GET("/message/notify", h.middleware.Handler, h.Notify)

	root.GET("/message/inbox", h.middleware.Handler, h.Inbox)
	root.DELETE("/message/inbox/:id", h.middleware.Handler, h.DeleteInbox)
	root.PUT("/message/inbox/:id", h.middleware.Handler, h.MarkInboxRead)
	root.POST("/message/inbox", h.middleware.Handler, h.PostInbox)
}

func (h *Message) PostInbox(ctx *gin.Context) {
	// extract the client from the context
	messageClient, ok := client.MessageFromContext(ctx)
	if !ok {
		log.Error(errors.ErrCaptchaNotCorrect)
		ctx.AbortWithError(http.StatusInternalServerError, errors.ErrClientNotFound)
		return
	}
}

func (h *Message) Inbox(ctx *gin.Context) {
	// extract the client from the context
	messageClient, ok := client.MessageFromContext(ctx)
	if !ok {
		log.Error(errors.ErrCaptchaNotCorrect)
		ctx.AbortWithError(http.StatusInternalServerError, errors.ErrClientNotFound)
		return
	}
}

func (h *Message) DeleteInbox(ctx *gin.Context) {
	messageClient, ok := client.MessageFromContext(ctx)
	if !ok {
		log.Error(errors.ErrCaptchaNotCorrect)
		ctx.AbortWithError(http.StatusInternalServerError, errors.ErrClientNotFound)
		return
	}
}

func (h *Message) MarkInboxRead(ctx *gin.Context) {
	messageClient, ok := client.MessageFromContext(ctx)
	if !ok {
		log.Error(errors.ErrCaptchaNotCorrect)
		ctx.AbortWithError(http.StatusInternalServerError, errors.ErrClientNotFound)
		return
	}
}

func (h *Message) Notify(ctx *gin.Context) {
	// TODO: Check acl
	conn, err := wsUpgrader.Upgrade(ctx.Writer, ctx.Request, nil)
	if err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	// Drop all incoming message
	go func() {
		_, _, err := conn.ReadMessage()
		if err != nil {
			conn.Close()
			return
		}
	}()

	// extract the client from the context
	messageClient, ok := client.MessageFromContext(ctx)
	if !ok {
		log.Error(errors.ErrCaptchaNotCorrect)
		ctx.AbortWithError(http.StatusInternalServerError, errors.ErrClientNotFound)
		return
	}

	notifyStream, err := messageClient.GetNotify(ctx, &proto.GetNotifyReq{})
	if err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	go func() {
		notifyItem, err := notifyStream.Recv()
		if err != nil {
			log.Error(err)
			notifyStream.Close()
			return
		}
		notifyJson, err := json.Marshal(notifyItem)
		if err != nil {
			log.Error(err)
			notifyStream.Close()
			return
		}
		conn.WriteMessage(websocket.TextMessage, notifyJson)
	}()
}
