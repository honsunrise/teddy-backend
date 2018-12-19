package handler

import (
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"github.com/zhsyourai/teddy-backend/api/clients"
	"github.com/zhsyourai/teddy-backend/common/errors"
	"github.com/zhsyourai/teddy-backend/common/proto"
	"net/http"
	"strconv"
	"strings"
)

func buildSort(sort string) ([]*proto.Sort, error) {
	rawSorts := strings.Split(sort, ",")
	sorts := make([]*proto.Sort, 0, len(rawSorts))
	for _, item := range rawSorts {
		nameAndOrder := strings.Split(item, ":")
		if len(nameAndOrder) == 1 {
			sorts = append(sorts, &proto.Sort{
				Name: nameAndOrder[0],
				Asc:  false,
			})
		} else if len(nameAndOrder) == 2 {
			if strings.ToUpper(nameAndOrder[1]) == "ASC" {
				sorts = append(sorts, &proto.Sort{
					Name: nameAndOrder[0],
					Asc:  true,
				})
			} else if strings.ToUpper(nameAndOrder[1]) == "DESC" {
				sorts = append(sorts, &proto.Sort{
					Name: nameAndOrder[0],
					Asc:  false,
				})
			} else {
				return nil, errors.ErrOrderNotCorrect
			}
		}
	}
	return sorts, nil
}

type Content struct {
}

func NewContentHandler() (*Content, error) {
	return &Content{}, nil
}

func (h *Content) HandlerNormal(root gin.IRoutes) {
	root.GET("/tags", h.GetAllTags)
	root.GET("/all", h.GetAllContents)
	root.GET("/search", h.Search)
	root.GET("/id/:id", h.GetContentDetail)
	root.POST("/id/:id/watch", h.WatchContent)
}

func (h *Content) HandlerAuth(root gin.IRoutes) {
	root.POST("/publish", h.PublishContent)
	root.POST("/id/:id", h.UpdateContent)
	root.DELETE("/id/:id", h.DeleteContent)

	root.GET("/favorite/user", h.GetUserFavThumb)
	root.GET("/favorite/id/:id", h.GetInfoFavThumb)
	root.POST("/favorite/id/:id", h.FavThumb)
	root.DELETE("/favorite/id/:id", h.DeleteFavoThumb)

	root.GET("/thumbUp/user", h.GetUserFavThumb)
	root.GET("/thumbUp/id/:id", h.GetInfoFavThumb)
	root.POST("/thumbUp/id/:id", h.FavThumb)
	root.DELETE("/thumbUp/id/:id", h.DeleteFavoThumb)

	root.GET("/thumbDown/user", h.GetUserFavThumb)
	root.GET("/thumbDown/id/:id", h.GetInfoFavThumb)
	root.POST("/thumbDown/id/:id", h.FavThumb)
	root.DELETE("/thumbDown/id/:id", h.DeleteFavoThumb)
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

func (h *Content) GetAllTags(ctx *gin.Context) {
	// extract the client from the context
	contentClient, ok := clients.ContentFromContext(ctx)
	if !ok {
		ctx.AbortWithError(http.StatusInternalServerError, errors.ErrClientNotFound)
		return
	}
	page, err := strconv.ParseUint(ctx.Query("page"), 10, 32)
	if err != nil && err.(*strconv.NumError).Num != "" {
		log.Error(errors.ErrCaptchaNotCorrect)
		ctx.AbortWithError(http.StatusInternalServerError, errors.ErrClientNotFound)
		return
	}

	size, err := strconv.ParseUint(ctx.Query("size"), 10, 32)
	if err != nil && err.(*strconv.NumError).Num != "" {
		log.Error(errors.ErrCaptchaNotCorrect)
		ctx.AbortWithError(http.StatusInternalServerError, errors.ErrClientNotFound)
		return
	}

	var sorts []*proto.Sort
	if ctx.Query("sort") != "" {
		sorts, err = buildSort(ctx.Query("sort"))
		if err != nil {
			log.Error(errors.ErrOrderNotCorrect)
			ctx.AbortWithError(http.StatusInternalServerError, errors.ErrOrderNotCorrect)
			return
		}
	}

	resp, err := contentClient.GetTags(ctx, &proto.GetTagReq{
		Page:  uint32(page),
		Size:  uint32(size),
		Sorts: sorts,
	})

	if err != nil {
		log.Error(err)
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	ctx.JSON(http.StatusOK, resp)
}

func (h *Content) GetAllContents(ctx *gin.Context) {
	// extract the client from the context
	contentClient, ok := clients.ContentFromContext(ctx)
	if !ok {
		ctx.AbortWithError(http.StatusInternalServerError, errors.ErrClientNotFound)
		return
	}
	page, err := strconv.ParseUint(ctx.Query("page"), 10, 32)
	if err != nil && err.(*strconv.NumError).Num != "" {
		log.Error(errors.ErrCaptchaNotCorrect)
		ctx.AbortWithError(http.StatusInternalServerError, errors.ErrClientNotFound)
		return
	}

	size, err := strconv.ParseUint(ctx.Query("size"), 10, 32)
	if err != nil && err.(*strconv.NumError).Num != "" {
		log.Error(errors.ErrCaptchaNotCorrect)
		ctx.AbortWithError(http.StatusInternalServerError, errors.ErrClientNotFound)
		return
	}

	var sorts []*proto.Sort
	if ctx.Query("sort") != "" {
		sorts, err = buildSort(ctx.Query("sort"))
		if err != nil {
			log.Error(errors.ErrOrderNotCorrect)
			ctx.AbortWithError(http.StatusInternalServerError, errors.ErrOrderNotCorrect)
			return
		}
	}

	resp, err := contentClient.GetInfos(ctx, &proto.GetInfosReq{
		Page:  uint32(page),
		Size:  uint32(size),
		Sorts: sorts,
	})

	if err != nil {
		log.Error(err)
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	ctx.JSON(http.StatusOK, resp)
}

func (h *Content) Search(ctx *gin.Context) {
	type okResp struct {
		Status string `json:"status"`
	}
	var jsonResp okResp
	jsonResp.Status = "OK"
	ctx.JSON(http.StatusOK, &jsonResp)
}

func (h *Content) PublishContent(ctx *gin.Context) {
	contentClient, ok := clients.ContentFromContext(ctx)
	if !ok {
		ctx.AbortWithError(http.StatusInternalServerError, errors.ErrClientNotFound)
		return
	}

	type publishInfoReq struct {
		Title          string            `json:"title"`
		Content        string            `json:"content"`
		Tags           []string          `json:"tags"`
		CanReview      bool              `json:"can_review"`
		CoverResources map[string]string `json:"cover_resources"`
	}
	var req publishInfoReq
	ctx.BindJSON(&req)

	_, err := contentClient.PublishInfo(ctx, &proto.PublishInfoReq{})

	if err != nil {
		log.Error(err)
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	ctx.Status(http.StatusOK)
}

func (h *Content) GetContentDetail(ctx *gin.Context) {
	contentClient, ok := clients.ContentFromContext(ctx)
	if !ok {
		ctx.AbortWithError(http.StatusInternalServerError, errors.ErrClientNotFound)
		return
	}

	infoID := ctx.Param("id")
	resp, err := contentClient.GetInfo(ctx, &proto.GetInfoReq{
		Id: infoID,
	})

	if err != nil {
		log.Error(err)
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	ctx.JSON(http.StatusOK, resp)
}

func (h *Content) UpdateContent(ctx *gin.Context) {
	contentClient, ok := clients.ContentFromContext(ctx)
	if !ok {
		ctx.AbortWithError(http.StatusInternalServerError, errors.ErrClientNotFound)
		return
	}

	type updateInfoReq struct {
		ID             string            `json:"id"`
		Title          string            `json:"title"`
		Content        string            `json:"content"`
		Tags           []string          `json:"tags"`
		CanReview      bool              `json:"can_review"`
		CoverResources map[string]string `json:"cover_resources"`
	}
	var req updateInfoReq
	ctx.BindJSON(&req)

	_, err := contentClient.EditInfo(ctx, &proto.EditInfoReq{})

	if err != nil {
		log.Error(err)
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	ctx.Status(http.StatusOK)
}

func (h *Content) DeleteContent(ctx *gin.Context) {
	contentClient, ok := clients.ContentFromContext(ctx)
	if !ok {
		ctx.AbortWithError(http.StatusInternalServerError, errors.ErrClientNotFound)
		return
	}

	infoID := ctx.Param("id")

	_, err := contentClient.DeleteInfo(ctx, &proto.InfoIDReq{
		InfoID: infoID,
	})

	if err != nil {
		log.Error(err)
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	ctx.Status(http.StatusOK)
}

func (h *Content) WatchContent(ctx *gin.Context) {
	contentClient, ok := clients.ContentFromContext(ctx)
	if !ok {
		ctx.AbortWithError(http.StatusInternalServerError, errors.ErrClientNotFound)
		return
	}

	infoID := ctx.Param("id")

	_, err := contentClient.WatchInfo(ctx, &proto.InfoIDReq{
		InfoID: infoID,
	})

	if err != nil {
		log.Error(err)
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	ctx.Status(http.StatusOK)
}

// For favorite
func (h *Content) GetUserFavThumb(ctx *gin.Context) {
	contentClient, ok := clients.ContentFromContext(ctx)
	if !ok {
		ctx.AbortWithError(http.StatusInternalServerError, errors.ErrClientNotFound)
		return
	}
	page, err := strconv.ParseUint(ctx.Query("page"), 10, 32)
	if err != nil && err.(*strconv.NumError).Num != "" {
		log.Error(errors.ErrCaptchaNotCorrect)
		ctx.AbortWithError(http.StatusInternalServerError, errors.ErrClientNotFound)
		return
	}

	size, err := strconv.ParseUint(ctx.Query("size"), 10, 32)
	if err != nil && err.(*strconv.NumError).Num != "" {
		log.Error(errors.ErrCaptchaNotCorrect)
		ctx.AbortWithError(http.StatusInternalServerError, errors.ErrClientNotFound)
		return
	}

	var sorts []*proto.Sort
	if ctx.Query("sort") != "" {
		sorts, err = buildSort(ctx.Query("sort"))
		if err != nil {
			log.Error(errors.ErrOrderNotCorrect)
			ctx.AbortWithError(http.StatusInternalServerError, errors.ErrOrderNotCorrect)
			return
		}
	}

	// TODO: get uid
	var uid string
	var resp *proto.InfoIDsResp
	if strings.Contains(ctx.Request.RequestURI, "favorite") {
		resp, err = contentClient.GetUserFavorite(ctx, &proto.UIDPageReq{
			Page:  uint32(page),
			Size:  uint32(size),
			Sorts: sorts,
			Uid:   uid,
		})
	} else if strings.Contains(ctx.Request.RequestURI, "thumbUp") {
		resp, err = contentClient.GetUserThumbUp(ctx, &proto.UIDPageReq{
			Page:  uint32(page),
			Size:  uint32(size),
			Sorts: sorts,
			Uid:   uid,
		})
	} else {
		resp, err = contentClient.GetUserThumbDown(ctx, &proto.UIDPageReq{
			Page:  uint32(page),
			Size:  uint32(size),
			Sorts: sorts,
			Uid:   uid,
		})
	}

	if err != nil {
		log.Error(err)
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	ctx.JSON(http.StatusOK, resp)
}

func (h *Content) GetInfoFavThumb(ctx *gin.Context) {
	contentClient, ok := clients.ContentFromContext(ctx)
	if !ok {
		ctx.AbortWithError(http.StatusInternalServerError, errors.ErrClientNotFound)
		return
	}
	page, err := strconv.ParseUint(ctx.Query("page"), 10, 32)
	if err != nil && err.(*strconv.NumError).Num != "" {
		log.Error(errors.ErrCaptchaNotCorrect)
		ctx.AbortWithError(http.StatusInternalServerError, errors.ErrClientNotFound)
		return
	}

	size, err := strconv.ParseUint(ctx.Query("size"), 10, 32)
	if err != nil && err.(*strconv.NumError).Num != "" {
		log.Error(errors.ErrCaptchaNotCorrect)
		ctx.AbortWithError(http.StatusInternalServerError, errors.ErrClientNotFound)
		return
	}

	var sorts []*proto.Sort
	if ctx.Query("sort") != "" {
		sorts, err = buildSort(ctx.Query("sort"))
		if err != nil {
			log.Error(errors.ErrOrderNotCorrect)
			ctx.AbortWithError(http.StatusInternalServerError, errors.ErrOrderNotCorrect)
			return
		}
	}

	infoID := ctx.Param("id")
	var resp *proto.UserIDsResp
	if strings.Contains(ctx.Request.RequestURI, "favorite") {
		resp, err = contentClient.GetInfoFavorite(ctx, &proto.InfoIDPageReq{
			Page:   uint32(page),
			Size:   uint32(size),
			Sorts:  sorts,
			InfoID: infoID,
		})
	} else if strings.Contains(ctx.Request.RequestURI, "thumbUp") {
		resp, err = contentClient.GetInfoThumbUp(ctx, &proto.InfoIDPageReq{
			Page:   uint32(page),
			Size:   uint32(size),
			Sorts:  sorts,
			InfoID: infoID,
		})
	} else {
		resp, err = contentClient.GetInfoThumbDown(ctx, &proto.InfoIDPageReq{
			Page:   uint32(page),
			Size:   uint32(size),
			Sorts:  sorts,
			InfoID: infoID,
		})
	}

	if err != nil {
		log.Error(err)
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	ctx.JSON(http.StatusOK, resp)
}

func (h *Content) FavThumb(ctx *gin.Context) {
	contentClient, ok := clients.ContentFromContext(ctx)
	if !ok {
		ctx.AbortWithError(http.StatusInternalServerError, errors.ErrClientNotFound)
		return
	}

	//TODO: fill uid
	var uid string
	var err error
	infoID := ctx.Param("id")
	if strings.Contains(ctx.Request.RequestURI, "favorite") {
		_, err = contentClient.Favorite(ctx, &proto.InfoIDAndUIDReq{
			InfoID: infoID,
			Uid:    uid,
		})
	} else if strings.Contains(ctx.Request.RequestURI, "thumbUp") {
		_, err = contentClient.ThumbUp(ctx, &proto.InfoIDAndUIDReq{
			InfoID: infoID,
			Uid:    uid,
		})
	} else {
		_, err = contentClient.ThumbDown(ctx, &proto.InfoIDAndUIDReq{
			InfoID: infoID,
			Uid:    uid,
		})
	}

	if err != nil {
		log.Error(err)
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	ctx.Status(http.StatusOK)
}

func (h *Content) DeleteFavoThumb(ctx *gin.Context) {
	contentClient, ok := clients.ContentFromContext(ctx)
	if !ok {
		ctx.AbortWithError(http.StatusInternalServerError, errors.ErrClientNotFound)
		return
	}

	//TODO: fill uid
	var uid string
	var err error
	infoID := ctx.Param("id")
	if strings.Contains(ctx.Request.RequestURI, "favorite") {
		_, err = contentClient.DeleteFavorite(ctx, &proto.InfoIDAndUIDReq{
			InfoID: infoID,
			Uid:    uid,
		})
	} else if strings.Contains(ctx.Request.RequestURI, "thumbUp") {
		_, err = contentClient.DeleteThumbUp(ctx, &proto.InfoIDAndUIDReq{
			InfoID: infoID,
			Uid:    uid,
		})
	} else {
		_, err = contentClient.DeleteThumbUp(ctx, &proto.InfoIDAndUIDReq{
			InfoID: infoID,
			Uid:    uid,
		})
	}

	if err != nil {
		log.Error(err)
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	ctx.Status(http.StatusOK)
}
