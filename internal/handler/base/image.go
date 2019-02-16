package base

import (
	"github.com/gin-gonic/gin"
	"github.com/minio/minio-go"
	"net/http"
	"teddy-backend/internal/gin_jwt"
	"teddy-backend/internal/handler/errors"
)

type Image struct {
	middleware  *gin_jwt.JwtMiddleware
	minioClient *minio.Client
	minioBucket string
}

func NewImageHandler(middleware *gin_jwt.JwtMiddleware, minioClient *minio.Client, bucket string) (*Image, error) {
	return &Image{
		middleware:  middleware,
		minioClient: minioClient,
		minioBucket: bucket,
	}, nil
}

func (h *Image) HandlerNormal(root gin.IRoutes) {
	root.GET("/:id", h.GetImage)
}

func (h *Image) HandlerAuth(root gin.IRoutes) {
}

func (h *Image) HandlerHealth(root gin.IRoutes) {
	root.Any("/", h.ReturnOK)
}

func (h *Image) ReturnOK(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{
		"status": "OK",
	})
}

func (h *Image) GetImage(ctx *gin.Context) {
	idStr := ctx.Param("id")

	ctx.Header("Cache-Control", "no-cache, no-store, must-revalidate")
	ctx.Header("Pragma", "no-cache")
	ctx.Header("Expires", "0")

	contentType := "application/octet-stream"
	download := ctx.Query("download") == "true"

	obj, err := h.minioClient.GetObject(h.minioBucket, idStr, minio.GetObjectOptions{})
	if err != nil {
		errors.AbortWithErrorJSON(ctx, errors.ErrUnknown)
		return
	}
	objStat, err := obj.Stat()
	if err != nil {
		errors.AbortWithErrorJSON(ctx, errors.ErrUnknown)
		return
	}
	if !download {
		contentType = objStat.ContentType
	}

	extraHeader := map[string]string{}
	ctx.DataFromReader(http.StatusOK, objStat.Size, contentType, obj, extraHeader)
}
