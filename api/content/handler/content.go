package handler

import (
	"bytes"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/golang/protobuf/jsonpb"
	"github.com/golang/protobuf/ptypes"
	"github.com/minio/minio-go"
	"github.com/zhsyourai/teddy-backend/api/clients"
	"github.com/zhsyourai/teddy-backend/common/gin_jwt"
	"github.com/zhsyourai/teddy-backend/common/proto/content"
	"net/http"
	"strings"
	"time"
)

const jsonContentType = "application/json; charset=utf-8"

type Content struct {
	middleware  *gin_jwt.JwtMiddleware
	minioClient *minio.Client
	minioBucket string
	marshaler   *jsonpb.Marshaler
}

func NewContentHandler(middleware *gin_jwt.JwtMiddleware,
	minioClient *minio.Client, bucket string) (*Content, error) {
	return &Content{
		middleware:  middleware,
		minioClient: minioClient,
		minioBucket: bucket,
		marshaler: &jsonpb.Marshaler{
			EnumsAsInts:  false,
			EmitDefaults: true,
		},
	}, nil
}

func (h *Content) HandlerNormal(root gin.IRoutes) {
	root.GET("/tags", h.GetTags)
	root.GET("/tags/:tagID", h.GetTag)

	root.GET("/info", h.GetInfos)
	root.GET("/info/:id", h.GetInfo)

	root.GET("/info/:id/segment", h.GetSegments)
	root.GET("/info/:id/segment/:segID", h.GetSegment)

	root.GET("/info/:id/segment/:segID/value", h.GetValues)
	root.GET("/info/:id/segment/:segID/value/:valID", h.GetValue)

	root.GET("/favorite/info/:id", h.GetInfoFavThumb)
	root.GET("/thumbUp/info/:id", h.GetInfoFavThumb)
	root.GET("/thumbDown/info/:id", h.GetInfoFavThumb)

	root.GET("/search", h.Search)
}

func (h *Content) HandlerAuth(root gin.IRoutes) {
	root.POST("/info", h.PublishInfo)
	root.POST("/info/:id", h.UpdateInfo)
	root.DELETE("/info/:id", h.DeleteInfo)

	root.POST("/info/:id/segment", h.PublishSegment)
	root.POST("/info/:id/segment/:segID", h.UpdateSegment)
	root.DELETE("/info/:id/segment/:segID", h.DeleteSegment)

	root.POST("/info/:id/segment/:segID/value", h.InsertValue)
	root.POST("/info/:id/segment/:segID/value/:valID", h.UpdateValue)
	root.DELETE("/info/:id/segment/:segID/value/:valID", h.DeleteValue)

	root.GET("/favorite/user", h.GetUserFavThumb)
	root.POST("/favorite/info/:id", h.FavThumb)
	root.DELETE("/favorite/info/:id", h.DeleteFavThumb)

	root.GET("/thumbUp/user", h.GetUserFavThumb)
	root.POST("/thumbUp/info/:id", h.FavThumb)
	root.DELETE("/thumbUp/info/:id", h.DeleteFavThumb)

	root.GET("/thumbDown/user", h.GetUserFavThumb)
	root.POST("/thumbDown/info/:id", h.FavThumb)
	root.DELETE("/thumbDown/info/:id", h.DeleteFavThumb)
}

func (h *Content) HandlerHealth(root gin.IRoutes) {
	root.Any("/", h.ReturnOK)
}

func (h *Content) ReturnOK(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{
		"status": "OK",
	})
}

func (h *Content) Search(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{
		"status": "OK",
	})
}

func (h *Content) GetTags(ctx *gin.Context) {
	contentClient := clients.ContentFromContext(ctx)

	page, size, sorts, err := extractPageSizeSort(ctx)
	if err != nil {
		ctx.Error(err)
		return
	}

	resp, err := contentClient.GetTags(ctx, &content.GetTagsReq{
		Type:  ctx.Query("type"),
		Page:  page,
		Size:  size,
		Sorts: sorts,
	})
	if err != nil {
		ctx.Error(err)
		return
	}

	var buf bytes.Buffer
	if err = h.marshaler.Marshal(&buf, resp); err != nil {
		ctx.Error(err)
		return
	}

	ctx.Data(http.StatusOK, jsonContentType, buf.Bytes())
}

func (h *Content) GetTag(ctx *gin.Context) {
	contentClient := clients.ContentFromContext(ctx)

	tagID := ctx.Param("tagID")
	resp, err := contentClient.GetTag(ctx, &content.GetTagReq{
		Id: tagID,
	})
	if err != nil {
		ctx.Error(err)
		return
	}

	var buf bytes.Buffer
	if err = h.marshaler.Marshal(&buf, resp); err != nil {
		ctx.Error(err)
		return
	}
	ctx.Data(http.StatusOK, jsonContentType, buf.Bytes())
}

func (h *Content) GetInfos(ctx *gin.Context) {
	contentClient := clients.ContentFromContext(ctx)
	principal := h.middleware.ExtractSub(ctx)

	page, size, sorts, err := extractPageSizeSort(ctx)
	if err != nil {
		ctx.Error(err)
		return
	}

	type queryReq struct {
		Tags             string    `form:"tags"`
		Country          string    `form:"country"`
		ContentTimeStart time.Time `form:"contentTimeStart" time_format:"2006-01-02T15:04:05Z07:00"`
		ContentTimeEnd   time.Time `form:"contentTimeEnd" time_format:"2006-01-02T15:04:05Z07:00"`
	}
	var req queryReq

	err = ctx.BindQuery(&req)
	if err != nil {
		ctx.Error(err)
		return
	}
	var tags []*content.TagAndType
	if req.Tags != "" {
		tags, err = buildTags(ctx.Query("tags"))
		if err != nil {
			ctx.Error(err)
			return
		}
	}

	contentTimeStart, err := ptypes.TimestampProto(req.ContentTimeStart)
	if err != nil {
		ctx.Error(err)
		return
	}
	contentTimeEnd, err := ptypes.TimestampProto(req.ContentTimeEnd)
	if err != nil {
		ctx.Error(err)
		return
	}

	resp, err := contentClient.GetInfos(ctx, &content.GetInfosReq{
		Uid:       principal,
		Page:      page,
		Size:      size,
		Tags:      tags,
		Sorts:     sorts,
		Country:   req.Country,
		StartTime: contentTimeStart,
		EndTime:   contentTimeEnd,
	})
	if err != nil {
		ctx.Error(err)
		return
	}

	var buf bytes.Buffer
	if err = h.marshaler.Marshal(&buf, resp); err != nil {
		ctx.Error(err)
		return
	}

	ctx.Data(http.StatusOK, jsonContentType, buf.Bytes())
}

func (h *Content) GetInfo(ctx *gin.Context) {
	contentClient := clients.ContentFromContext(ctx)
	principal := h.middleware.ExtractSub(ctx)

	infoID := ctx.Param("id")
	resp, err := contentClient.GetInfo(ctx, &content.GetInfoReq{
		InfoID: infoID,
		Uid:    principal,
	})
	if err != nil {
		ctx.Error(err)
		return
	}

	var buf bytes.Buffer
	if err = h.marshaler.Marshal(&buf, resp); err != nil {
		ctx.Error(err)
		return
	}
	ctx.Data(http.StatusOK, jsonContentType, buf.Bytes())
}

func (h *Content) PublishInfo(ctx *gin.Context) {
	contentClient := clients.ContentFromContext(ctx)
	principal := h.middleware.ExtractSub(ctx)

	type publishInfoReq struct {
		Title       string    `form:"title" binding:"required"`
		Author      string    `form:"author" binding:"required"`
		Summary     string    `form:"summary" binding:"required"`
		Country     string    `form:"country" binding:"required"`
		Tags        []string  `form:"tags" binding:"required"`
		CanReview   bool      `form:"canReview" binding:"required"`
		ContentTime time.Time `form:"contentTime" binding:"required" time_format:"2006-01-02T15:04:05Z07:00"`
	}
	var req publishInfoReq
	err := ctx.Bind(&req)
	if err != nil {
		return
	}

	contentTime, err := ptypes.TimestampProto(req.ContentTime)
	if err != nil {
		ctx.Error(err)
		return
	}

	form, err := ctx.MultipartForm()
	if err != nil {
		ctx.String(http.StatusBadRequest, fmt.Sprintf("get form err: %s", err.Error()))
		return
	}
	covers := form.File["covers"]
	if len(covers) == 0 {
		ctx.String(http.StatusBadRequest, fmt.Sprintf("at lest have one cover"))
		return
	}

	coverResources := make(map[string]string)
	for i, cover := range covers {
		objectName, err := uploadFile(cover, h.minioClient, h.minioBucket)
		if err != nil {
			ctx.String(http.StatusBadRequest, fmt.Sprintf("upload file err: %s", err.Error()))
			return
		}
		coverResources[fmt.Sprintf("%d", i)] = objectName
	}

	tags := make([]*content.TagAndType, 0, len(req.Tags))
	for _, v := range req.Tags {
		arr := strings.Split(v, ":")
		tags = append(tags, &content.TagAndType{
			Type: arr[0],
			Tag:  arr[1],
		})
	}

	resp, err := contentClient.PublishInfo(ctx, &content.PublishInfoReq{
		Uid:            principal,
		Author:         req.Author,
		Summary:        req.Summary,
		Title:          req.Title,
		Country:        req.Country,
		Tags:           tags,
		CanReview:      req.CanReview,
		CoverResources: coverResources,
		ContentTime:    contentTime,
	})

	if err != nil {
		ctx.Error(err)
		return
	}

	var buf bytes.Buffer
	if err = h.marshaler.Marshal(&buf, resp); err != nil {
		ctx.Error(err)
		return
	}
	ctx.Data(http.StatusOK, jsonContentType, buf.Bytes())
}

func (h *Content) UpdateInfo(ctx *gin.Context) {
	contentClient := clients.ContentFromContext(ctx)
	principal := h.middleware.ExtractSub(ctx)

	type updateInfoReq struct {
		Title       string    `form:"title" binding:"required"`
		Author      string    `form:"author" binding:"required"`
		Summary     string    `form:"summary" binding:"required"`
		Country     string    `form:"country" binding:"required"`
		Tags        []string  `form:"tags" binding:"required"`
		CanReview   bool      `form:"canReview" binding:"required"`
		ContentTime time.Time `form:"contentTime" binding:"required" time_format:"2006-01-02T15:04:05Z07:00"`
	}
	var req updateInfoReq
	err := ctx.Bind(&req)
	if err != nil {
		return
	}

	contentTime, err := ptypes.TimestampProto(req.ContentTime)
	if err != nil {
		ctx.Error(err)
		return
	}

	form, err := ctx.MultipartForm()
	if err != nil {
		ctx.Error(err)
		return
	}
	covers := form.File["covers"]

	coverResources := make(map[string]string)
	for i, cover := range covers {
		objectName, err := uploadFile(cover, h.minioClient, h.minioBucket)
		if err != nil {
			ctx.String(http.StatusBadRequest, fmt.Sprintf("upload file err: %s", err.Error()))
			return
		}
		coverResources[fmt.Sprintf("%d", i)] = objectName
	}

	tags := make([]*content.TagAndType, 0, len(req.Tags))
	for _, v := range req.Tags {
		arr := strings.Split(v, ":")
		tags = append(tags, &content.TagAndType{
			Type: arr[0],
			Tag:  arr[1],
		})
	}

	infoID := ctx.Param("id")
	_, err = contentClient.EditInfo(ctx, &content.EditInfoReq{
		Uid:            principal,
		InfoID:         infoID,
		Title:          req.Title,
		Author:         req.Author,
		Country:        req.Country,
		Summary:        req.Summary,
		Tags:           tags,
		CanReview:      req.CanReview,
		CoverResources: coverResources,
		ContentTime:    contentTime,
	})

	if err != nil {
		ctx.Error(err)
		return
	}

	ctx.Status(http.StatusOK)
}

func (h *Content) DeleteInfo(ctx *gin.Context) {
	contentClient := clients.ContentFromContext(ctx)
	principal := h.middleware.ExtractSub(ctx)

	infoID := ctx.Param("id")
	_, err := contentClient.GetInfo(ctx, &content.GetInfoReq{
		Uid:    principal,
		InfoID: infoID,
	})
	if err != nil {
		ctx.Error(err)
		return
	}

	_, err = contentClient.DeleteInfo(ctx, &content.InfoOneReq{
		InfoID: infoID,
	})

	if err != nil {
		ctx.Error(err)
		return
	}

	ctx.Status(http.StatusOK)
}

// For favorite
func (h *Content) GetUserFavThumb(ctx *gin.Context) {
	contentClient := clients.ContentFromContext(ctx)

	principal := h.middleware.ExtractSub(ctx)

	page, size, sorts, err := extractPageSizeSort(ctx)
	if err != nil {
		ctx.Error(err)
		return
	}

	var resp *content.InfoIDsResp
	if strings.Contains(ctx.Request.RequestURI, "favorite") {
		resp, err = contentClient.GetUserFavorite(ctx, &content.UIDPageReq{
			Page:  page,
			Size:  size,
			Sorts: sorts,
			Uid:   principal,
		})
	} else if strings.Contains(ctx.Request.RequestURI, "thumbUp") {
		resp, err = contentClient.GetUserThumbUp(ctx, &content.UIDPageReq{
			Page:  page,
			Size:  size,
			Sorts: sorts,
			Uid:   principal,
		})
	} else {
		resp, err = contentClient.GetUserThumbDown(ctx, &content.UIDPageReq{
			Page:  page,
			Size:  size,
			Sorts: sorts,
			Uid:   principal,
		})
	}

	var buf bytes.Buffer
	if err = h.marshaler.Marshal(&buf, resp); err != nil {
		ctx.Error(err)
		return
	}
	ctx.Data(http.StatusOK, jsonContentType, buf.Bytes())
}

func (h *Content) GetInfoFavThumb(ctx *gin.Context) {
	contentClient := clients.ContentFromContext(ctx)

	page, size, sorts, err := extractPageSizeSort(ctx)
	if err != nil {
		ctx.Error(err)
		return
	}

	infoID := ctx.Param("id")
	var resp *content.UserIDsResp
	if strings.Contains(ctx.Request.RequestURI, "favorite") {
		resp, err = contentClient.GetInfoFavorite(ctx, &content.InfoIDPageReq{
			Page:   page,
			Size:   size,
			Sorts:  sorts,
			InfoID: infoID,
		})
	} else if strings.Contains(ctx.Request.RequestURI, "thumbUp") {
		resp, err = contentClient.GetInfoThumbUp(ctx, &content.InfoIDPageReq{
			Page:   page,
			Size:   size,
			Sorts:  sorts,
			InfoID: infoID,
		})
	} else {
		resp, err = contentClient.GetInfoThumbDown(ctx, &content.InfoIDPageReq{
			Page:   page,
			Size:   size,
			Sorts:  sorts,
			InfoID: infoID,
		})
	}

	if err != nil {
		ctx.Error(err)
		return
	}

	var buf bytes.Buffer
	if err = h.marshaler.Marshal(&buf, resp); err != nil {
		ctx.Error(err)
		return
	}
	ctx.Data(http.StatusOK, jsonContentType, buf.Bytes())
}

func (h *Content) FavThumb(ctx *gin.Context) {
	contentClient := clients.ContentFromContext(ctx)
	principal := h.middleware.ExtractSub(ctx)

	var err error
	infoID := ctx.Param("id")
	if strings.Contains(ctx.Request.RequestURI, "favorite") {
		_, err = contentClient.Favorite(ctx, &content.InfoIDWithUIDReq{
			InfoID: infoID,
			Uid:    principal,
		})
	} else if strings.Contains(ctx.Request.RequestURI, "thumbUp") {
		_, err = contentClient.ThumbUp(ctx, &content.InfoIDWithUIDReq{
			InfoID: infoID,
			Uid:    principal,
		})
	} else {
		_, err = contentClient.ThumbDown(ctx, &content.InfoIDWithUIDReq{
			InfoID: infoID,
			Uid:    principal,
		})
	}

	if err != nil {
		ctx.Error(err)
		return
	}

	ctx.Status(http.StatusOK)
}

func (h *Content) DeleteFavThumb(ctx *gin.Context) {
	contentClient := clients.ContentFromContext(ctx)
	principal := h.middleware.ExtractSub(ctx)

	var err error
	infoID := ctx.Param("id")
	if strings.Contains(ctx.Request.RequestURI, "favorite") {
		_, err = contentClient.DeleteFavorite(ctx, &content.InfoIDWithUIDReq{
			InfoID: infoID,
			Uid:    principal,
		})
	} else if strings.Contains(ctx.Request.RequestURI, "thumbUp") {
		_, err = contentClient.DeleteThumbUp(ctx, &content.InfoIDWithUIDReq{
			InfoID: infoID,
			Uid:    principal,
		})
	} else {
		_, err = contentClient.DeleteThumbDown(ctx, &content.InfoIDWithUIDReq{
			InfoID: infoID,
			Uid:    principal,
		})
	}

	if err != nil {
		ctx.Error(err)
		return
	}

	ctx.Status(http.StatusOK)
}

func (h *Content) GetSegments(ctx *gin.Context) {
	contentClient := clients.ContentFromContext(ctx)

	page, size, sorts, err := extractPageSizeSort(ctx)
	if err != nil {
		ctx.Error(err)
		return
	}

	var labels []string
	if ctx.Query("labels") != "" {
		labels = strings.Split(ctx.Query("labels"), ",")
	}

	infoID := ctx.Param("id")
	resp, err := contentClient.GetSegments(ctx, &content.GetSegmentsReq{
		InfoID: infoID,
		Labels: labels,
		Page:   page,
		Size:   size,
		Sorts:  sorts,
	})

	if err != nil {
		ctx.Error(err)
		return
	}

	var buf bytes.Buffer
	if err = h.marshaler.Marshal(&buf, resp); err != nil {
		ctx.Error(err)
		return
	}
	ctx.Data(http.StatusOK, jsonContentType, buf.Bytes())
}

func (h *Content) GetSegment(ctx *gin.Context) {
	contentClient := clients.ContentFromContext(ctx)

	infoID := ctx.Param("id")
	segID := ctx.Param("segID")

	resp, err := contentClient.GetSegment(ctx, &content.SegmentOneReq{
		InfoID: infoID,
		SegID:  segID,
	})

	if err != nil {
		ctx.Error(err)
		return
	}

	var buf bytes.Buffer
	if err = h.marshaler.Marshal(&buf, resp); err != nil {
		ctx.Error(err)
		return
	}
	ctx.Data(http.StatusOK, jsonContentType, buf.Bytes())
}

func (h *Content) PublishSegment(ctx *gin.Context) {
	contentClient := clients.ContentFromContext(ctx)
	principal := h.middleware.ExtractSub(ctx)

	infoID := ctx.Param("id")
	_, err := contentClient.GetInfo(ctx, &content.GetInfoReq{
		Uid:    principal,
		InfoID: infoID,
	})

	if err != nil {
		ctx.Error(err)
		return
	}

	type publishSegmentReq struct {
		No     uint64   `form:"no" binding:"required"`
		Labels []string `form:"labels" binding:"required"`
		Title  string   `form:"title" binding:"required"`
	}
	var req publishSegmentReq
	err = ctx.Bind(&req)
	if err != nil {
		ctx.Error(err)
		return
	}

	resp, err := contentClient.PublishSegment(ctx, &content.PublishSegmentReq{
		InfoID: infoID,
		No:     req.No,
		Labels: req.Labels,
		Title:  req.Title,
	})

	if err != nil {
		ctx.Error(err)
		return
	}

	var buf bytes.Buffer
	if err = h.marshaler.Marshal(&buf, resp); err != nil {
		ctx.Error(err)
		return
	}
	ctx.Data(http.StatusOK, jsonContentType, buf.Bytes())
}

func (h *Content) UpdateSegment(ctx *gin.Context) {
	contentClient := clients.ContentFromContext(ctx)
	principal := h.middleware.ExtractSub(ctx)

	infoID := ctx.Param("id")
	segID := ctx.Param("segID")

	type updateSegmentReq struct {
		No     uint64   `json:"no"`
		Labels []string `json:"labels"`
	}
	var req updateSegmentReq
	err := ctx.Bind(&req)
	if err != nil {
		return
	}

	_, err = contentClient.GetInfo(ctx, &content.GetInfoReq{
		Uid:    principal,
		InfoID: infoID,
	})

	if err != nil {
		ctx.Error(err)
		return
	}

	_, err = contentClient.EditSegment(ctx, &content.EditSegmentReq{
		InfoID:  infoID,
		SegID:   segID,
		No:      req.No,
		Labels:  req.Labels,
		Content: nil,
	})

	if err != nil {
		ctx.Error(err)
		return
	}

	ctx.Status(http.StatusOK)
}

func (h *Content) DeleteSegment(ctx *gin.Context) {
	contentClient := clients.ContentFromContext(ctx)
	principal := h.middleware.ExtractSub(ctx)

	infoID := ctx.Param("id")
	segID := ctx.Param("segID")

	_, err := contentClient.GetInfo(ctx, &content.GetInfoReq{
		Uid:    principal,
		InfoID: infoID,
	})

	if err != nil {
		ctx.Error(err)
		return
	}

	_, err = contentClient.DeleteSegment(ctx, &content.SegmentOneReq{
		InfoID: infoID,
		SegID:  segID,
	})

	if err != nil {
		ctx.Error(err)
		return
	}

	ctx.Status(http.StatusOK)
}

func (h *Content) GetValues(ctx *gin.Context) {
	contentClient := clients.ContentFromContext(ctx)

	infoID := ctx.Param("id")
	segID := ctx.Param("segID")

	page, size, sorts, err := extractPageSizeSort(ctx)
	if err != nil {
		ctx.Error(err)
		return
	}

	resp, err := contentClient.GetValues(ctx, &content.GetValuesReq{
		InfoID: infoID,
		SegID:  segID,
		Page:   page,
		Size:   size,
		Sorts:  sorts,
	})

	if err != nil {
		ctx.Error(err)
		return
	}

	var buf bytes.Buffer
	if err = h.marshaler.Marshal(&buf, resp); err != nil {
		ctx.Error(err)
		return
	}
	ctx.Data(http.StatusOK, jsonContentType, buf.Bytes())
}

func (h *Content) GetValue(ctx *gin.Context) {
	contentClient := clients.ContentFromContext(ctx)

	infoID := ctx.Param("id")
	segID := ctx.Param("segID")
	valID := ctx.Param("valID")

	resp, err := contentClient.GetValue(ctx, &content.ValueOneReq{
		InfoID: infoID,
		SegID:  segID,
		ValID:  valID,
	})
	if err != nil {
		ctx.Error(err)
		return
	}

	var buf bytes.Buffer
	if err = h.marshaler.Marshal(&buf, resp); err != nil {
		ctx.Error(err)
		return
	}
	ctx.Data(http.StatusOK, jsonContentType, buf.Bytes())
}

func (h *Content) InsertValue(ctx *gin.Context) {
	contentClient := clients.ContentFromContext(ctx)
	principal := h.middleware.ExtractSub(ctx)

	infoID := ctx.Param("id")
	segID := ctx.Param("segID")

	_, err := contentClient.GetInfo(ctx, &content.GetInfoReq{
		Uid:    principal,
		InfoID: infoID,
	})

	if err != nil {
		ctx.Error(err)
		return
	}

	form, err := ctx.MultipartForm()
	if err != nil {
		ctx.String(http.StatusBadRequest, fmt.Sprintf("get form err: %s", err.Error()))
		return
	}
	file := form.File["file"]
	if len(file) > 1 || len(file) == 0 {
		ctx.String(http.StatusBadRequest, fmt.Sprintf("get form err"))
		return
	}

	objectName, err := uploadFile(file[0], h.minioClient, h.minioBucket)
	if err != nil {
		ctx.String(http.StatusBadRequest, fmt.Sprintf("upload file err: %s", err.Error()))
		return
	}

	resp, err := contentClient.InsertValue(ctx, &content.InsertValueReq{
		InfoID: infoID,
		SegID:  segID,
		Time:   ptypes.TimestampNow(),
		Value:  objectName,
	})

	if err != nil {
		ctx.Error(err)
		return
	}

	var buf bytes.Buffer
	if err = h.marshaler.Marshal(&buf, resp); err != nil {
		ctx.Error(err)
		return
	}
	ctx.Data(http.StatusOK, jsonContentType, buf.Bytes())
}

func (h *Content) UpdateValue(ctx *gin.Context) {
	contentClient := clients.ContentFromContext(ctx)
	principal := h.middleware.ExtractSub(ctx)

	infoID := ctx.Param("id")
	segID := ctx.Param("segID")
	valID := ctx.Param("valID")

	_, err := contentClient.GetInfo(ctx, &content.GetInfoReq{
		Uid:    principal,
		InfoID: infoID,
	})

	if err != nil {
		ctx.Error(err)
		return
	}

	form, err := ctx.MultipartForm()
	if err != nil {
		ctx.String(http.StatusBadRequest, fmt.Sprintf("get form err: %s", err.Error()))
		return
	}
	file := form.File["file"]
	if len(file) > 1 || len(file) == 0 {
		ctx.String(http.StatusBadRequest, fmt.Sprintf("get form err"))
		return
	}

	objectName, err := uploadFile(file[0], h.minioClient, h.minioBucket)
	if err != nil {
		ctx.String(http.StatusBadRequest, fmt.Sprintf("upload file err: %s", err.Error()))
		return
	}

	_, err = contentClient.EditValue(ctx, &content.EditValueReq{
		InfoID: infoID,
		SegID:  segID,
		ValID:  valID,
		Time:   ptypes.TimestampNow(),
		Value:  objectName,
	})

	if err != nil {
		ctx.Error(err)
		return
	}

	ctx.Status(http.StatusOK)
}

func (h *Content) DeleteValue(ctx *gin.Context) {
	contentClient := clients.ContentFromContext(ctx)
	principal := h.middleware.ExtractSub(ctx)

	infoID := ctx.Param("id")
	segID := ctx.Param("segID")
	valID := ctx.Param("valID")

	_, err := contentClient.GetInfo(ctx, &content.GetInfoReq{
		Uid:    principal,
		InfoID: infoID,
	})

	if err != nil {
		ctx.Error(err)
		return
	}

	_, err = contentClient.DeleteValue(ctx, &content.ValueOneReq{
		InfoID: infoID,
		SegID:  segID,
		ValID:  valID,
	})

	if err != nil {
		ctx.Error(err)
		return
	}

	ctx.Status(http.StatusOK)
}
