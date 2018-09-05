package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/prometheus/common/log"
	"github.com/zhsyourai/teddy-backend/api/base/client"
	"github.com/zhsyourai/teddy-backend/api/gin-jwt"
	"github.com/zhsyourai/teddy-backend/captcha/proto"
	"github.com/zhsyourai/teddy-backend/common/errors"
	"net/http"
	"path"
	"strings"
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
	root.GET("/base/captcha/:id", h.GetCaptchaImage)

	root.PUT("/base/profile/:id")
	root.GET("/base/profile/:id")
	root.POST("/base/profile/:id/avatar")
}

func (h *Base) GetCaptchaImage(ctx *gin.Context) {
	// extract the client from the context
	captchaClient, ok := client.CaptchaFromContext(ctx)
	if !ok {
		log.Error(errors.ErrCaptchaNotCorrect)
		ctx.AbortWithError(http.StatusInternalServerError, errors.ErrClientNotFound)
		return
	}

	idStr := ctx.Param("id")
	ext := path.Ext(idStr)
	id := idStr[:len(idStr)-len(ext)]
	if ext == "" || id == "" {
		ctx.Status(http.StatusNotFound)
		return
	}

	lang := strings.ToLower(ctx.Param("lang"))
	// Fill header
	if ctx.Param("download") != "true" {
		ctx.Header("Content-Type", "application/octet-stream")
	}
	ctx.Header("Cache-Control", "no-cache, no-store, must-revalidate")
	ctx.Header("Pragma", "no-cache")
	ctx.Header("Expires", "0")

	switch ext {
	case ".png":
		resp, err := captchaClient.GetImage(ctx, &proto.GetImageReq{
			Len:    6,
			Width:  240,
			Height: 80,
			Reload: ctx.Param("reload") != "",
		})
		if err != nil {
			log.Error(err)
			ctx.AbortWithError(http.StatusInternalServerError, err)
			return
		}
		ctx.Data(http.StatusOK, "image/png", resp.Image)
	case ".wav":
		resp, err := captchaClient.GetVoice(ctx, &proto.GetVoiceReq{
			Len:    6,
			Lang:   lang,
			Reload: ctx.Param("reload") != "",
		})
		if err != nil {
			log.Error(err)
			ctx.AbortWithError(http.StatusInternalServerError, err)
			return
		}
		ctx.Data(http.StatusOK, "audio/x-wav", resp.VoiceWav)
	default:
		ctx.Status(http.StatusNotFound)
		return
	}
}
