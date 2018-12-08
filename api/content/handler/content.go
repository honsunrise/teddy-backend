package handler

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

type Content struct {
}

func NewContentHandler() (*Content, error) {
	return &Content{}, nil
}

func (h *Content) HandlerNormal(root gin.IRoutes) {
}

func (h *Content) HandlerHealth(root gin.IRoutes) {
	root.Any("/", h.ReturnOK)
}

func (h *Content) ReturnOK(ctx *gin.Context) {
	type okResp struct {
		Status string `json:"status"`
	}
	var jsonResp okResp
	jsonResp.Status = "OK"
	ctx.JSON(http.StatusOK, &jsonResp)
}
