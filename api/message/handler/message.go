package handler

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/prometheus/common/log"
	"github.com/zhsyourai/teddy-backend/api/clients"
	"github.com/zhsyourai/teddy-backend/common/proto/message"
	"github.com/zhsyourai/teddy-backend/common/types"
	"net/http"
	"time"
)

type Message struct {
}

func NewMessageHandler() (*Message, error) {
	return &Message{}, nil
}

var wsUpgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func (h *Message) HandlerNormal(root gin.IRoutes) {
	root.GET("/system")
}

func (h *Message) HandlerAuth(root gin.IRoutes) {
	root.GET("/notify", h.Notify)
	root.GET("/inbox", h.Inbox)
	root.DELETE("/inbox/:id", h.DeleteInbox)
	root.PUT("/inbox/:id", h.MarkInboxRead)
	root.POST("/inbox", h.PostInbox)
}

func (h *Message) HandlerHealth(root gin.IRoutes) {
	root.Any("/", h.ReturnOK)
}

func (h *Message) ReturnOK(ctx *gin.Context) {
	type okResp struct {
		Status string `json:"status"`
	}
	var jsonResp okResp
	jsonResp.Status = "OK"
	ctx.JSON(http.StatusOK, &jsonResp)
}

func (h *Message) PostInbox(ctx *gin.Context) {
	// extract the client from the context
	_, ok := clients.MessageFromContext(ctx)
	if !ok {
		log.Error(ErrCaptchaNotCorrect)
		ctx.AbortWithError(http.StatusInternalServerError, ErrClientNotFound)
		return
	}
}

func (h *Message) Inbox(ctx *gin.Context) {
	// extract the client from the context
	paging := &types.Paging{}
	if err := ctx.BindQuery(paging); err != nil {
		log.Error(err)
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	uidStr := ctx.Query("uid")

	messageClient, ok := clients.MessageFromContext(ctx)
	if !ok {
		log.Error(ErrCaptchaNotCorrect)
		ctx.AbortWithError(http.StatusInternalServerError, ErrClientNotFound)
		return
	}
	inboxResp, err := messageClient.GetInBox(ctx, &message.GetInBoxReq{
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
	_, ok := clients.MessageFromContext(ctx)
	if !ok {
		log.Error(ErrCaptchaNotCorrect)
		ctx.AbortWithError(http.StatusInternalServerError, ErrClientNotFound)
		return
	}
}

func (h *Message) MarkInboxRead(ctx *gin.Context) {
	_, ok := clients.MessageFromContext(ctx)
	if !ok {
		log.Error(ErrCaptchaNotCorrect)
		ctx.AbortWithError(http.StatusInternalServerError, ErrClientNotFound)
		return
	}
}

func (h *Message) Notify(ctx *gin.Context) {
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
	messageClient, ok := clients.MessageFromContext(ctx)
	if !ok {
		log.Error(ErrCaptchaNotCorrect)
		ctx.AbortWithError(http.StatusInternalServerError, ErrClientNotFound)
		return
	}

	notifyStream, err := messageClient.GetNotify(ctx, &message.GetNotifyReq{})
	if err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	go func() {
		notifyItem, err := notifyStream.Recv()
		if err != nil {
			log.Error(err)
			return
		}
		notifyJson, err := json.Marshal(notifyItem)
		if err != nil {
			log.Error(err)
			return
		}
		conn.WriteMessage(websocket.TextMessage, notifyJson)
	}()
}
