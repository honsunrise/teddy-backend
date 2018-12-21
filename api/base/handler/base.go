package handler

import (
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"github.com/zhsyourai/teddy-backend/api/clients"
	"github.com/zhsyourai/teddy-backend/common/proto/captcha"
	"net/http"
	"path"
	"strings"
)

type Base struct {
}

func NewBaseHandler() (*Base, error) {
	return &Base{}, nil
}

func (h *Base) HandlerNormal(root gin.IRoutes) {
	root.GET("/captcha", h.GetCaptchaId)
	root.GET("/captcha/:id", h.GetCaptchaData)

	root.GET("/profile/:id")
}

func (h *Base) HandlerAuth(root gin.IRoutes) {
	root.PUT("/profile/:id")
	root.GET("/profile/:id/detail")
	root.POST("/profile/:id/avatar")
}

func (h *Base) HandlerHealth(root gin.IRoutes) {
	root.Any("/", h.ReturnOK)
}

func (h *Base) ReturnOK(ctx *gin.Context) {
	type okResp struct {
		Status string `json:"status"`
	}
	var jsonResp okResp
	jsonResp.Status = "OK"
	ctx.JSON(http.StatusOK, &jsonResp)
}

func (h *Base) GetCaptchaId(ctx *gin.Context) {
	// extract the client from the context
	captchaClient, ok := clients.CaptchaFromContext(ctx)
	if !ok {
		log.Error(ErrClientNotFound)
		ctx.AbortWithError(http.StatusInternalServerError, ErrClientNotFound)
		return
	}

	idResp, err := captchaClient.GetCaptchaId(ctx, &captcha.GetCaptchaIdReq{
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
	captchaClient, ok := clients.CaptchaFromContext(ctx)
	if !ok {
		log.Error(ErrClientNotFound)
		ctx.AbortWithError(http.StatusInternalServerError, ErrClientNotFound)
		return
	}

	idStr := ctx.Param("id")
	ext := path.Ext(idStr)
	id := idStr[:len(idStr)-len(ext)]
	if ext == "" || id == "" {
		ctx.Status(http.StatusNotFound)
		return
	}

	// Fill header
	ctx.Header("Cache-Control", "no-cache, no-store, must-revalidate")
	ctx.Header("Pragma", "no-cache")
	ctx.Header("Expires", "0")

	contentType := "application/octet-stream"
	reload := ctx.Query("reload") == "true"
	download := ctx.Query("download") == "true"
	switch ext {
	case ".png":
		resp, err := captchaClient.GetImageData(ctx, &captcha.GetImageDataReq{
			Id:     id,
			Width:  280,
			Height: 93,
			Reload: reload,
		})
		if err != nil {
			log.Error(err)
			ctx.AbortWithError(http.StatusInternalServerError, err)
			return
		}
		if !download {
			contentType = "image/png"
		}
		ctx.Data(http.StatusOK, contentType, resp.Image)
	case ".wav":
		lang := strings.ToLower(ctx.Query("lang"))
		resp, err := captchaClient.GetVoiceData(ctx, &captcha.GetVoiceDataReq{
			Id:     id,
			Lang:   lang,
			Reload: reload,
		})
		if err != nil {
			log.Error(err)
			ctx.AbortWithError(http.StatusInternalServerError, err)
			return
		}
		if !download {
			contentType = "audio/x-wav"
		}
		ctx.Data(http.StatusOK, contentType, resp.VoiceWav)
	default:
		ctx.Status(http.StatusNotFound)
		return
	}
}
