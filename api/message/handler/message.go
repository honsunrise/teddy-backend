package handler

import (
	"encoding/json"
	"github.com/casbin/casbin"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/prometheus/common/log"
	"github.com/zhsyourai/teddy-backend/api/gin-jwt"
	"github.com/zhsyourai/teddy-backend/api/message/client"
	"github.com/zhsyourai/teddy-backend/common/errors"
	"github.com/zhsyourai/teddy-backend/common/types"
	"github.com/zhsyourai/teddy-backend/message/proto"
	"gopkg.in/dgrijalva/jwt-go.v3"
	"net/http"
	"time"
)

type Message struct {
	enforcer   *casbin.Enforcer
	middleware *gin_jwt.JwtMiddleware
}

func NewMessageHandler(enforcer *casbin.Enforcer, middleware *gin_jwt.JwtMiddleware) (*Message, error) {
	return &Message{
		enforcer:   enforcer,
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
	_, ok := client.MessageFromContext(ctx)
	if !ok {
		log.Error(errors.ErrCaptchaNotCorrect)
		ctx.AbortWithError(http.StatusInternalServerError, errors.ErrClientNotFound)
		return
	}
}

func (h *Message) Inbox(ctx *gin.Context) {
	token, err := h.middleware.ExtractToken(ctx)
	if err != nil {
		ctx.AbortWithStatus(http.StatusUnauthorized)
	}

	sub := token.Claims.(jwt.MapClaims)["uid"] // the user that wants to access a resource.
	obj := "uaa.changePassword"                // the resource that is going to be accessed.
	act := "read,write"                        // the operation that the user performs on the resource.

	if h.enforcer.Enforce(sub, obj, act) != true {
		ctx.AbortWithStatus(http.StatusForbidden)
		return
	}

	// extract the client from the context
	paging := &types.Paging{}
	if err := ctx.BindQuery(paging); err != nil {
		log.Error(err)
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	uidStr := ctx.Query("uid")

	messageClient, ok := client.MessageFromContext(ctx)
	if !ok {
		log.Error(errors.ErrCaptchaNotCorrect)
		ctx.AbortWithError(http.StatusInternalServerError, errors.ErrClientNotFound)
		return
	}
	inboxResp, err := messageClient.GetInBox(ctx, &proto.GetInBoxReq{
		Page: paging.Page,
		Size: paging.Size,
		Uid:  uidStr,
	})
	if err != nil {
		log.Error(err)
		ctx.AbortWithError(http.StatusInternalServerError, err)
	}
	type inboxItem struct {
		Id       string    `json:"id,omitempty"`
		Topic    string    `json:"topic,omitempty"`
		Content  string    `json:"content,omitempty"`
		From     string    `json:"from,omitempty"`
		Type     uint32    `json:"type,omitempty"`
		Unread   bool      `json:"unread,omitempty"`
		SendTime time.Time `json:"sendTime,omitempty"`
		ReadTime time.Time `json:"ReadTime,omitempty"`
	}

	result := make([]inboxItem, len(inboxResp.Items))
	for i, item := range inboxResp.Items {
		result[i].Id = item.Id
		result[i].Topic = item.Topic
		result[i].Content = item.Content
		result[i].From = item.From
		result[i].Type = item.Type
		result[i].Unread = item.Unread
		result[i].SendTime = time.Unix(item.SendTime.Seconds, int64(item.SendTime.Nanos))
		result[i].ReadTime = time.Unix(item.ReadTime.Seconds, int64(item.ReadTime.Nanos))
	}
	ctx.JSON(http.StatusOK, result)
}

func (h *Message) DeleteInbox(ctx *gin.Context) {
	_, ok := client.MessageFromContext(ctx)
	if !ok {
		log.Error(errors.ErrCaptchaNotCorrect)
		ctx.AbortWithError(http.StatusInternalServerError, errors.ErrClientNotFound)
		return
	}
}

func (h *Message) MarkInboxRead(ctx *gin.Context) {
	_, ok := client.MessageFromContext(ctx)
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
