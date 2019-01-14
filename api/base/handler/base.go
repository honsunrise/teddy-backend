package handler

import (
	"github.com/gin-gonic/gin"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"net/http"
	"path"
	"strings"
	"teddy-backend/api/clients"
	"teddy-backend/api/errors"
	"teddy-backend/common/gin_jwt"
	"teddy-backend/common/proto/captcha"
)

type Base struct {
	middleware *gin_jwt.JwtMiddleware
}

func NewBaseHandler(middleware *gin_jwt.JwtMiddleware) (*Base, error) {
	return &Base{
		middleware: middleware,
	}, nil
}

func (h *Base) HandlerNormal(root gin.IRoutes) {
	root.GET("/captcha", h.GetCaptchaId)
	root.GET("/captcha/:id", h.GetCaptchaData)

	root.GET("/profile/:id")
}

func (h *Base) HandlerAuth(root gin.IRoutes) {
	root.POST("/profile/:id")
}

func (h *Base) HandlerHealth(root gin.IRoutes) {
	root.Any("/", h.ReturnOK)
}

func (h *Base) ReturnOK(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{
		"status": "OK",
	})
}

func (h *Base) GetCaptchaId(ctx *gin.Context) {
	captchaClient := clients.CaptchaFromContext(ctx)

	idResp, err := captchaClient.GetCaptchaId(ctx, &captcha.GetCaptchaIdReq{
		Len: 6,
	})
	if err != nil {
		errors.AbortWithErrorJSON(ctx, errors.ErrUnknown)
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
	captchaClient := clients.CaptchaFromContext(ctx)

	idStr := ctx.Param("id")
	ext := path.Ext(idStr)
	id := idStr[:len(idStr)-len(ext)]
	if id == "" {
		errors.AbortWithErrorJSON(ctx, errors.ErrCaptchaIDNotFound)
		return
	}
	if ext == "" {
		errors.AbortWithErrorJSON(ctx, errors.ErrCaptchaExtNotSupport)
		return
	}

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
			if status.Code(err) == codes.NotFound {
				errors.AbortWithErrorJSON(ctx, errors.ErrCaptchaIDNotFound)
			} else {
				errors.AbortWithErrorJSON(ctx, errors.ErrUnknown)
			}
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
			if status.Code(err) == codes.NotFound {
				errors.AbortWithErrorJSON(ctx, errors.ErrCaptchaIDNotFound)
			} else {
				errors.AbortWithErrorJSON(ctx, errors.ErrUnknown)
			}
			return
		}
		if !download {
			contentType = "audio/x-wav"
		}
		ctx.Data(http.StatusOK, contentType, resp.VoiceWav)
	default:
		errors.AbortWithErrorJSON(ctx, errors.ErrCaptchaExtNotSupport)
		return
	}
}
