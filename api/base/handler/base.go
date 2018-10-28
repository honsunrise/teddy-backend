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

func NewBaseHandler() (*Base, error) {
	return &Base{}, nil
}

func (h *Base) Handler(root gin.IRoutes) {
	root.GET("/base/captcha", h.GetCaptchaId)
	root.GET("/base/captcha/:id", h.GetCaptchaData)

	root.PUT("/base/profile/:id")
	root.GET("/base/profile/:id")
	root.POST("/base/profile/:id/avatar")
}

func (h *Base) GetCaptchaId(ctx *gin.Context) {
	// extract the client from the context
	captchaClient, ok := client.CaptchaFromContext(ctx)
	if !ok {
		log.Error(errors.ErrCaptchaNotCorrect)
		ctx.AbortWithError(http.StatusInternalServerError, errors.ErrClientNotFound)
		return
	}

	idResp, err := captchaClient.GetCaptchaId(ctx, &proto.GetCaptchaIdReq{
		Len: 6,
	})
	if err != nil {
		log.Error(err)
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	type captchaIdResp struct {
		Id string `json:"id"`
	}
	var jsonResp captchaIdResp
	jsonResp.Id = idResp.Id
	ctx.JSON(http.StatusOK, &jsonResp)
}

func (h *Base) GetCaptchaData(ctx *gin.Context) {
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
	ctx.Header("Cache-Control", "no-cache, no-store, must-revalidate")
	ctx.Header("Pragma", "no-cache")
	ctx.Header("Expires", "0")

	contentType := "application/octet-stream"

	switch ext {
	case ".png":
		resp, err := captchaClient.GetImageData(ctx, &proto.GetImageDataReq{
			Id:     id,
			Width:  240,
			Height: 80,
			Reload: ctx.Param("reload") != "",
		})
		if err != nil {
			log.Error(err)
			ctx.AbortWithError(http.StatusInternalServerError, err)
			return
		}
		if ctx.Param("download") != "true" {
			contentType = "image/png"
		}
		ctx.Data(http.StatusOK, contentType, resp.Image)
	case ".wav":
		resp, err := captchaClient.GetVoiceData(ctx, &proto.GetVoiceDataReq{
			Id:     id,
			Lang:   lang,
			Reload: ctx.Param("reload") != "",
		})
		if err != nil {
			log.Error(err)
			ctx.AbortWithError(http.StatusInternalServerError, err)
			return
		}
		if ctx.Param("download") != "true" {
			contentType = "audio/x-wav"
		}
		ctx.Data(http.StatusOK, contentType, resp.VoiceWav)
	default:
		ctx.Status(http.StatusNotFound)
		return
	}
}
