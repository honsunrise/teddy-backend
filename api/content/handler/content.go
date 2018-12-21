package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/golang/protobuf/ptypes"
	log "github.com/sirupsen/logrus"
	"github.com/zhsyourai/teddy-backend/api/clients"
	"github.com/zhsyourai/teddy-backend/common/proto/content"
	"net/http"
	"strconv"
	"strings"
	"time"
)

func buildSort(sort string) ([]*content.Sort, error) {
	rawSorts := strings.Split(sort, ",")
	sorts := make([]*content.Sort, 0, len(rawSorts))
	for _, item := range rawSorts {
		nameAndOrder := strings.Split(item, ":")
		if len(nameAndOrder) == 1 {
			sorts = append(sorts, &content.Sort{
				Name: nameAndOrder[0],
				Asc:  false,
			})
		} else if len(nameAndOrder) == 2 {
			if strings.ToUpper(nameAndOrder[1]) == "ASC" {
				sorts = append(sorts, &content.Sort{
					Name: nameAndOrder[0],
					Asc:  true,
				})
			} else if strings.ToUpper(nameAndOrder[1]) == "DESC" {
				sorts = append(sorts, &content.Sort{
					Name: nameAndOrder[0],
					Asc:  false,
				})
			} else {
				return nil, ErrOrderNotCorrect
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
		ctx.AbortWithError(http.StatusInternalServerError, ErrClientNotFound)
		return
	}
	page, err := strconv.ParseUint(ctx.Query("page"), 10, 32)
	if err != nil && err.(*strconv.NumError).Num != "" {
		log.Error(ErrCaptchaNotCorrect)
		ctx.AbortWithError(http.StatusInternalServerError, ErrClientNotFound)
		return
	}

	size, err := strconv.ParseUint(ctx.Query("size"), 10, 32)
	if err != nil && err.(*strconv.NumError).Num != "" {
		log.Error(ErrCaptchaNotCorrect)
		ctx.AbortWithError(http.StatusInternalServerError, ErrClientNotFound)
		return
	}

	var sorts []*content.Sort
	if ctx.Query("sort") != "" {
		sorts, err = buildSort(ctx.Query("sort"))
		if err != nil {
			log.Error(ErrOrderNotCorrect)
			ctx.AbortWithError(http.StatusInternalServerError, ErrOrderNotCorrect)
			return
		}
	}

	resp, err := contentClient.GetTags(ctx, &content.GetTagReq{
		Page:  uint32(page),
		Size:  uint32(size),
		Sorts: sorts,
	})

	if err != nil {
		log.Error(err)
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	type tagsResult struct {
		Tag         string    `json:"tag"`
		Usage       uint64    `json:"usage"`
		CreateTime  time.Time `json:"create_time"`
		LastUseTime time.Time `json:"last_use_time"`
	}

	results := make([]*tagsResult, 0, len(resp.Tags))
	for _, tag := range resp.Tags {
		createTime, err := ptypes.Timestamp(tag.CreateTime)
		if err != nil {
			log.Error(err)
			continue
		}

		lastUseTime, err := ptypes.Timestamp(tag.LastUseTime)
		if err != nil {
			log.Error(err)
			continue
		}
		results = append(results, &tagsResult{
			Tag:         tag.Tag,
			Usage:       tag.Usage,
			CreateTime:  createTime,
			LastUseTime: lastUseTime,
		})
	}

	ctx.JSON(http.StatusOK, results)
}

func (h *Content) GetAllContents(ctx *gin.Context) {
	// extract the client from the context
	contentClient, ok := clients.ContentFromContext(ctx)
	if !ok {
		ctx.AbortWithError(http.StatusInternalServerError, ErrClientNotFound)
		return
	}
	page, err := strconv.ParseUint(ctx.Query("page"), 10, 32)
	if err != nil && err.(*strconv.NumError).Num != "" {
		log.Error(ErrCaptchaNotCorrect)
		ctx.AbortWithError(http.StatusInternalServerError, ErrClientNotFound)
		return
	}

	size, err := strconv.ParseUint(ctx.Query("size"), 10, 32)
	if err != nil && err.(*strconv.NumError).Num != "" {
		log.Error(ErrCaptchaNotCorrect)
		ctx.AbortWithError(http.StatusInternalServerError, ErrClientNotFound)
		return
	}

	var sorts []*content.Sort
	if ctx.Query("sort") != "" {
		sorts, err = buildSort(ctx.Query("sort"))
		if err != nil {
			log.Error(ErrOrderNotCorrect)
			ctx.AbortWithError(http.StatusInternalServerError, ErrOrderNotCorrect)
			return
		}
	}

	resp, err := contentClient.GetInfos(ctx, &content.GetInfosReq{
		Page:  uint32(page),
		Size:  uint32(size),
		Sorts: sorts,
	})

	if err != nil {
		log.Error(err)
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	type infosResult struct {
		Id             string            `json:"id"`
		UID            string            `json:"uid"`
		Title          string            `json:"title"`
		Content        string            `json:"content"`
		ContentTime    time.Time         `json:"contentTime"`
		CoverResources map[string]string `json:"coverResources"`
		PublishTime    time.Time         `json:"publishTime"`
		LastReviewTime time.Time         `json:"lastReviewTime"`
		Valid          bool              `json:"valid"`
		WatchCount     int64             `json:"watchCount"`
		Tags           []string          `json:"tags"`
		ThumbUp        int64             `json:"thumbUp"`
		IsThumbUp      bool              `json:"isThumbUp"`
		ThumbUpList    []string          `json:"thumbUpList"`
		ThumbDown      int64             `json:"thumbDown"`
		IsThumbDown    bool              `json:"isThumbDown"`
		ThumbDownList  []string          `json:"thumbDownList"`
		Favorites      int64             `json:"favorites"`
		IsFavorite     bool              `json:"isFavorite"`
		FavoriteList   []string          `json:"favoriteList"`
		LastModifyTime time.Time         `json:"lastModifyTime"`
		CanReview      bool              `json:"canReview"`
	}

	results := make([]*infosResult, 0, len(resp.Infos))
	for _, info := range resp.Infos {
		publishTime, err := ptypes.Timestamp(info.PublishTime)
		if err != nil {
			log.Error(err)
		}

		lastReviewTime, err := ptypes.Timestamp(info.LastReviewTime)
		if err != nil {
			log.Error(err)
		}

		lastModifyTime, err := ptypes.Timestamp(info.LastModifyTime)
		if err != nil {
			log.Error(err)
		}

		contentTime, err := ptypes.Timestamp(info.ContentTime)
		if err != nil {
			log.Error(err)
		}

		results = append(results, &infosResult{
			Id:             info.InfoID,
			UID:            info.Uid,
			Title:          info.Title,
			Content:        info.Content,
			ContentTime:    contentTime,
			CoverResources: info.CoverResources,
			PublishTime:    publishTime,
			LastReviewTime: lastReviewTime,
			Valid:          info.Valid,
			WatchCount:     info.WatchCount,
			Tags:           info.Tags,
			ThumbUp:        info.ThumbUp,
			IsThumbUp:      info.IsThumbUp,
			ThumbUpList:    info.ThumbUpList,
			ThumbDown:      info.ThumbDown,
			IsThumbDown:    info.IsThumbDown,
			ThumbDownList:  info.ThumbDownList,
			Favorites:      info.Favorites,
			IsFavorite:     info.IsFavorite,
			FavoriteList:   info.FavoriteList,
			LastModifyTime: lastModifyTime,
			CanReview:      info.CanReview,
		})
	}

	ctx.JSON(http.StatusOK, results)
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
		ctx.AbortWithError(http.StatusInternalServerError, ErrClientNotFound)
		return
	}

	type publishInfoReq struct {
		Title          string            `json:"title"`
		Content        string            `json:"content"`
		Tags           []string          `json:"tags"`
		CanReview      bool              `json:"can_review"`
		CoverResources map[string]string `json:"cover_resources"`
		ContentTime    time.Time         `json:"content_time"`
	}
	var req publishInfoReq
	ctx.BindJSON(&req)

	//TODO: Get uid
	uid := "7791850604"
	_, err := contentClient.PublishInfo(ctx, &content.PublishInfoReq{
		Uid:            uid,
		Title:          req.Title,
		Content:        req.Content,
		Tags:           req.Tags,
		CanReview:      req.CanReview,
		CoverResources: req.CoverResources,
	})

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
		ctx.AbortWithError(http.StatusInternalServerError, ErrClientNotFound)
		return
	}

	infoID := ctx.Param("id")
	info, err := contentClient.GetInfo(ctx, &content.GetInfoReq{
		InfoID: infoID,
		Uid:    "7791850604",
	})

	if err != nil {
		log.Error(err)
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	type infosResult struct {
		Id             string            `json:"id"`
		UID            string            `json:"uid"`
		Title          string            `json:"title"`
		Content        string            `json:"content"`
		ContentTime    time.Time         `json:"contentTime"`
		CoverResources map[string]string `json:"coverResources"`
		PublishTime    time.Time         `json:"publishTime"`
		LastReviewTime time.Time         `json:"lastReviewTime"`
		Valid          bool              `json:"valid"`
		WatchCount     int64             `json:"watchCount"`
		Tags           []string          `json:"tags"`
		ThumbUp        int64             `json:"thumbUp"`
		IsThumbUp      bool              `json:"isThumbUp"`
		ThumbUpList    []string          `json:"thumbUpList"`
		ThumbDown      int64             `json:"thumbDown"`
		IsThumbDown    bool              `json:"isThumbDown"`
		ThumbDownList  []string          `json:"thumbDownList"`
		Favorites      int64             `json:"favorites"`
		IsFavorite     bool              `json:"isFavorite"`
		FavoriteList   []string          `json:"favoriteList"`
		LastModifyTime time.Time         `json:"lastModifyTime"`
		CanReview      bool              `json:"canReview"`
	}

	publishTime, err := ptypes.Timestamp(info.PublishTime)
	if err != nil {
		log.Error(err)
	}

	lastReviewTime, err := ptypes.Timestamp(info.LastReviewTime)
	if err != nil {
		log.Error(err)
	}

	lastModifyTime, err := ptypes.Timestamp(info.LastModifyTime)
	if err != nil {
		log.Error(err)
	}

	contentTime, err := ptypes.Timestamp(info.ContentTime)
	if err != nil {
		log.Error(err)
	}

	resp := &infosResult{
		Id:             info.InfoID,
		UID:            info.Uid,
		Title:          info.Title,
		Content:        info.Content,
		ContentTime:    contentTime,
		CoverResources: info.CoverResources,
		PublishTime:    publishTime,
		LastReviewTime: lastReviewTime,
		Valid:          info.Valid,
		WatchCount:     info.WatchCount,
		Tags:           info.Tags,
		ThumbUp:        info.ThumbUp,
		IsThumbUp:      info.IsThumbUp,
		ThumbUpList:    info.ThumbUpList,
		ThumbDown:      info.ThumbDown,
		IsThumbDown:    info.IsThumbDown,
		ThumbDownList:  info.ThumbDownList,
		Favorites:      info.Favorites,
		IsFavorite:     info.IsFavorite,
		FavoriteList:   info.FavoriteList,
		LastModifyTime: lastModifyTime,
		CanReview:      info.CanReview,
	}
	ctx.JSON(http.StatusOK, resp)
}

func (h *Content) UpdateContent(ctx *gin.Context) {
	contentClient, ok := clients.ContentFromContext(ctx)
	if !ok {
		ctx.AbortWithError(http.StatusInternalServerError, ErrClientNotFound)
		return
	}

	type updateInfoReq struct {
		Title          string            `json:"title"`
		Content        string            `json:"content"`
		Tags           []string          `json:"tags"`
		CanReview      bool              `json:"can_review"`
		CoverResources map[string]string `json:"cover_resources"`
	}
	var req updateInfoReq
	ctx.BindJSON(&req)

	infoID := ctx.Param("id")
	_, err := contentClient.EditInfo(ctx, &content.EditInfoReq{
		InfoID:         infoID,
		Title:          req.Title,
		Content:        req.Content,
		Tags:           req.Tags,
		CanReview:      req.CanReview,
		CoverResources: req.CoverResources,
	})

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
		ctx.AbortWithError(http.StatusInternalServerError, ErrClientNotFound)
		return
	}

	infoID := ctx.Param("id")

	_, err := contentClient.DeleteInfo(ctx, &content.InfoIDReq{
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
		ctx.AbortWithError(http.StatusInternalServerError, ErrClientNotFound)
		return
	}

	infoID := ctx.Param("id")

	_, err := contentClient.WatchInfo(ctx, &content.InfoIDReq{
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
		ctx.AbortWithError(http.StatusInternalServerError, ErrClientNotFound)
		return
	}
	page, err := strconv.ParseUint(ctx.Query("page"), 10, 32)
	if err != nil && err.(*strconv.NumError).Num != "" {
		log.Error(ErrCaptchaNotCorrect)
		ctx.AbortWithError(http.StatusInternalServerError, ErrClientNotFound)
		return
	}

	size, err := strconv.ParseUint(ctx.Query("size"), 10, 32)
	if err != nil && err.(*strconv.NumError).Num != "" {
		log.Error(ErrCaptchaNotCorrect)
		ctx.AbortWithError(http.StatusInternalServerError, ErrClientNotFound)
		return
	}

	var sorts []*content.Sort
	if ctx.Query("sort") != "" {
		sorts, err = buildSort(ctx.Query("sort"))
		if err != nil {
			log.Error(ErrOrderNotCorrect)
			ctx.AbortWithError(http.StatusInternalServerError, ErrOrderNotCorrect)
			return
		}
	}

	// TODO: get uid
	var uid string
	var resp *content.InfoIDsResp
	if strings.Contains(ctx.Request.RequestURI, "favorite") {
		resp, err = contentClient.GetUserFavorite(ctx, &content.UIDPageReq{
			Page:  uint32(page),
			Size:  uint32(size),
			Sorts: sorts,
			Uid:   uid,
		})
	} else if strings.Contains(ctx.Request.RequestURI, "thumbUp") {
		resp, err = contentClient.GetUserThumbUp(ctx, &content.UIDPageReq{
			Page:  uint32(page),
			Size:  uint32(size),
			Sorts: sorts,
			Uid:   uid,
		})
	} else {
		resp, err = contentClient.GetUserThumbDown(ctx, &content.UIDPageReq{
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
		ctx.AbortWithError(http.StatusInternalServerError, ErrClientNotFound)
		return
	}
	page, err := strconv.ParseUint(ctx.Query("page"), 10, 32)
	if err != nil && err.(*strconv.NumError).Num != "" {
		log.Error(ErrCaptchaNotCorrect)
		ctx.AbortWithError(http.StatusInternalServerError, ErrClientNotFound)
		return
	}

	size, err := strconv.ParseUint(ctx.Query("size"), 10, 32)
	if err != nil && err.(*strconv.NumError).Num != "" {
		log.Error(ErrCaptchaNotCorrect)
		ctx.AbortWithError(http.StatusInternalServerError, ErrClientNotFound)
		return
	}

	var sorts []*content.Sort
	if ctx.Query("sort") != "" {
		sorts, err = buildSort(ctx.Query("sort"))
		if err != nil {
			log.Error(ErrOrderNotCorrect)
			ctx.AbortWithError(http.StatusInternalServerError, ErrOrderNotCorrect)
			return
		}
	}

	infoID := ctx.Param("id")
	var resp *content.UserIDsResp
	if strings.Contains(ctx.Request.RequestURI, "favorite") {
		resp, err = contentClient.GetInfoFavorite(ctx, &content.InfoIDPageReq{
			Page:   uint32(page),
			Size:   uint32(size),
			Sorts:  sorts,
			InfoID: infoID,
		})
	} else if strings.Contains(ctx.Request.RequestURI, "thumbUp") {
		resp, err = contentClient.GetInfoThumbUp(ctx, &content.InfoIDPageReq{
			Page:   uint32(page),
			Size:   uint32(size),
			Sorts:  sorts,
			InfoID: infoID,
		})
	} else {
		resp, err = contentClient.GetInfoThumbDown(ctx, &content.InfoIDPageReq{
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
		ctx.AbortWithError(http.StatusInternalServerError, ErrClientNotFound)
		return
	}

	//TODO: fill uid
	uid := "7791850604"
	var err error
	infoID := ctx.Param("id")
	if strings.Contains(ctx.Request.RequestURI, "favorite") {
		_, err = contentClient.Favorite(ctx, &content.InfoIDAndUIDReq{
			InfoID: infoID,
			Uid:    uid,
		})
	} else if strings.Contains(ctx.Request.RequestURI, "thumbUp") {
		_, err = contentClient.ThumbUp(ctx, &content.InfoIDAndUIDReq{
			InfoID: infoID,
			Uid:    uid,
		})
	} else {
		_, err = contentClient.ThumbDown(ctx, &content.InfoIDAndUIDReq{
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
		ctx.AbortWithError(http.StatusInternalServerError, ErrClientNotFound)
		return
	}

	//TODO: fill uid
	var uid string
	var err error
	infoID := ctx.Param("id")
	if strings.Contains(ctx.Request.RequestURI, "favorite") {
		_, err = contentClient.DeleteFavorite(ctx, &content.InfoIDAndUIDReq{
			InfoID: infoID,
			Uid:    uid,
		})
	} else if strings.Contains(ctx.Request.RequestURI, "thumbUp") {
		_, err = contentClient.DeleteThumbUp(ctx, &content.InfoIDAndUIDReq{
			InfoID: infoID,
			Uid:    uid,
		})
	} else {
		_, err = contentClient.DeleteThumbUp(ctx, &content.InfoIDAndUIDReq{
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
