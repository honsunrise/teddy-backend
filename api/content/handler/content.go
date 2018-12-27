package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/golang/protobuf/ptypes"
	log "github.com/sirupsen/logrus"
	"github.com/zhsyourai/teddy-backend/api/clients"
	"github.com/zhsyourai/teddy-backend/api/gin_jwt"
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

func buildTags(tag string) ([]*content.TagAndType, error) {
	rawTags := strings.Split(tag, ",")
	tags := make([]*content.TagAndType, 0, len(rawTags))
	for _, item := range rawTags {
		typeAndTag := strings.Split(item, ":")
		if len(typeAndTag) == 2 {
			tags = append(tags, &content.TagAndType{
				Type: typeAndTag[0],
				Tag:  typeAndTag[1],
			})

		} else {
			return nil, ErrTagNotCorrect
		}
	}
	return tags, nil
}

type Content struct {
	middleware *gin_jwt.JwtMiddleware
}

func NewContentHandler(middleware *gin_jwt.JwtMiddleware) (*Content, error) {
	return &Content{
		middleware: middleware,
	}, nil
}

func (h *Content) HandlerNormal(root gin.IRoutes) {
	root.GET("/tags", h.GetAllTags)

	root.GET("/info", h.GetAllInfos)
	root.GET("/info/:id", h.GetInfoDetail)

	root.GET("/info/:id/segment", h.GetAllSegments)
	root.GET("/info/:id/segment/:segID", h.GetSegmentDetail)

	root.GET("/search", h.Search)
}

func (h *Content) HandlerAuth(root gin.IRoutes) {
	root.POST("/info", h.PublishInfo)
	root.POST("/info/:id", h.UpdateInfo)
	root.DELETE("/info/:id", h.DeleteInfo)

	root.POST("/info/:id/segment", h.PublishSegment)
	root.POST("/info/:id/segment/:segID", h.UpdateSegment)
	root.DELETE("/info/:id/segment/:segID", h.DeleteSegment)

	root.GET("/favorite/user", h.GetUserFavThumb)
	root.GET("/favorite/info/:id", h.GetInfoFavThumb)
	root.POST("/favorite/info/:id", h.FavThumb)
	root.DELETE("/favorite/info/:id", h.DeleteFavThumb)

	root.GET("/thumbUp/user", h.GetUserFavThumb)
	root.GET("/thumbUp/info/:id", h.GetInfoFavThumb)
	root.POST("/thumbUp/info/:id", h.FavThumb)
	root.DELETE("/thumbUp/info/:id", h.DeleteFavThumb)

	root.GET("/thumbDown/user", h.GetUserFavThumb)
	root.GET("/thumbDown/info/:id", h.GetInfoFavThumb)
	root.POST("/thumbDown/info/:id", h.FavThumb)
	root.DELETE("/thumbDown/info/:id", h.DeleteFavThumb)
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

func (h *Content) Search(ctx *gin.Context) {
	type okResp struct {
		Status string `json:"status"`
	}
	var jsonResp okResp
	jsonResp.Status = "OK"
	ctx.JSON(http.StatusOK, &jsonResp)
}

func (h *Content) GetAllTags(ctx *gin.Context) {
	principal := ""
	authPayload, err := h.middleware.ExtractToken(ctx)
	if err != gin_jwt.ErrContextNotHaveToken {
		principal = authPayload["sub"].(string)
	}
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
		Uid:   principal,
		Type:  ctx.Query("type"),
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
		Type        string    `json:"type"`
		Tag         string    `json:"tag"`
		Usage       uint64    `json:"usage"`
		CreateTime  time.Time `json:"createTime"`
		LastUseTime time.Time `json:"lastUseTime"`
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
			Type:        tag.Type,
			Tag:         tag.Tag,
			Usage:       tag.Usage,
			CreateTime:  createTime,
			LastUseTime: lastUseTime,
		})
	}

	ctx.JSON(http.StatusOK, gin.H{
		"total_count": resp.TotalCount,
		"items":       results,
	})
}

func (h *Content) GetAllInfos(ctx *gin.Context) {
	principal := ""
	authPayload, err := h.middleware.ExtractToken(ctx)
	if err != gin_jwt.ErrContextNotHaveToken {
		principal = authPayload["sub"].(string)
	}
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
	if ctx.Query("sorts") != "" {
		sorts, err = buildSort(ctx.Query("sorts"))
		if err != nil {
			ctx.AbortWithError(http.StatusInternalServerError, err)
			return
		}
	}

	var tags []*content.TagAndType
	if ctx.Query("tags") != "" {
		tags, err = buildTags(ctx.Query("tags"))
		if err != nil {
			ctx.AbortWithError(http.StatusInternalServerError, err)
			return
		}
	}
	resp, err := contentClient.GetInfos(ctx, &content.GetInfosReq{
		Uid:   principal,
		Page:  uint32(page),
		Size:  uint32(size),
		Tags:  tags,
		Sorts: sorts,
	})

	if err != nil {
		log.Error(err)
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	type infoResultTag struct {
		Type string `json:"type"`
		Tag  string `json:"tag"`
	}

	type infoResult struct {
		Id              string            `json:"id"`
		UID             string            `json:"uid"`
		Author          string            `json:"author"`
		Title           string            `json:"title"`
		Summary         string            `json:"summary"`
		Country         string            `json:"country"`
		ContentTime     time.Time         `json:"contentTime"`
		CoverResources  map[string]string `json:"coverResources"`
		PublishTime     time.Time         `json:"publishTime"`
		LastReviewTime  time.Time         `json:"lastReviewTime"`
		Valid           bool              `json:"valid"`
		WatchCount      uint64            `json:"watchCount"`
		Tags            []*infoResultTag  `json:"tags"`
		ThumbUp         uint64            `json:"thumbUps"`
		IsThumbUp       bool              `json:"isThumbUp"`
		ThumbUpList     []string          `json:"thumbUpList"`
		ThumbDown       uint64            `json:"thumbDowns"`
		IsThumbDown     bool              `json:"isThumbDown"`
		ThumbDownList   []string          `json:"thumbDownList"`
		Favorites       uint64            `json:"favorites"`
		IsFavorite      bool              `json:"isFavorite"`
		FavoriteList    []string          `json:"favoriteList"`
		LastModifyTime  time.Time         `json:"lastModifyTime"`
		CanReview       bool              `json:"canReview"`
		Archived        bool              `json:"archived"`
		LatestSegmentID string            `json:"latestSegmentID"`
		SegmentCount    uint64            `json:"segmentCount"`
	}
	results := make([]*infoResult, 0, len(resp.Infos))
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

		tags := make([]*infoResultTag, 0, len(info.Tags))
		for _, v := range info.Tags {
			tags = append(tags, &infoResultTag{
				Type: v.Type,
				Tag:  v.Tag,
			})
		}

		results = append(results, &infoResult{
			Id:              info.InfoID,
			UID:             info.Uid,
			Title:           info.Title,
			Author:          info.Author,
			Summary:         info.Summary,
			Country:         info.Country,
			ContentTime:     contentTime,
			CoverResources:  info.CoverResources,
			PublishTime:     publishTime,
			LastReviewTime:  lastReviewTime,
			Valid:           info.Valid,
			WatchCount:      info.WatchCount,
			Tags:            tags,
			ThumbUp:         info.ThumbUps,
			IsThumbUp:       info.IsThumbUp,
			ThumbUpList:     info.ThumbUpList,
			ThumbDown:       info.ThumbDowns,
			IsThumbDown:     info.IsThumbDown,
			ThumbDownList:   info.ThumbDownList,
			Favorites:       info.Favorites,
			IsFavorite:      info.IsFavorite,
			FavoriteList:    info.FavoriteList,
			LastModifyTime:  lastModifyTime,
			CanReview:       info.CanReview,
			Archived:        info.Archived,
			LatestSegmentID: info.LatestSegmentID,
			SegmentCount:    info.SegmentCount,
		})
	}

	ctx.JSON(http.StatusOK, gin.H{
		"total_count": resp.TotalCount,
		"items":       results,
	})
}

func (h *Content) PublishInfo(ctx *gin.Context) {
	authPayload, err := h.middleware.ExtractToken(ctx)
	if err == gin_jwt.ErrContextNotHaveToken {
		ctx.Error(err)
		return
	}
	principal := authPayload["sub"].(string)

	contentClient, ok := clients.ContentFromContext(ctx)
	if !ok {
		ctx.AbortWithError(http.StatusInternalServerError, ErrClientNotFound)
		return
	}

	type publishInfoTag struct {
		Type string `json:"type"`
		Tag  string `json:"tag"`
	}

	type publishInfoReq struct {
		Title          string            `json:"title"`
		Author         string            `json:"author"`
		Summary        string            `json:"summary"`
		Country        string            `json:"country"`
		Tags           []*publishInfoTag `json:"tags"`
		CanReview      bool              `json:"canReview"`
		CoverResources map[string]string `json:"coverResources"`
		ContentTime    time.Time         `json:"contentTime"`
	}
	var req publishInfoReq
	ctx.Bind(&req)

	tags := make([]*content.TagAndType, 0, len(req.Tags))
	for _, v := range req.Tags {
		tags = append(tags, &content.TagAndType{
			Type: v.Type,
			Tag:  v.Tag,
		})
	}

	_, err = contentClient.PublishInfo(ctx, &content.PublishInfoReq{
		Uid:            principal,
		Author:         req.Author,
		Summary:        req.Summary,
		Title:          req.Title,
		Country:        req.Country,
		Tags:           tags,
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

func (h *Content) GetInfoDetail(ctx *gin.Context) {
	principal := ""
	authPayload, err := h.middleware.ExtractToken(ctx)
	if err != gin_jwt.ErrContextNotHaveToken {
		principal = authPayload["sub"].(string)
	}

	contentClient, ok := clients.ContentFromContext(ctx)
	if !ok {
		ctx.AbortWithError(http.StatusInternalServerError, ErrClientNotFound)
		return
	}

	infoID := ctx.Param("id")
	info, err := contentClient.GetInfo(ctx, &content.GetInfoReq{
		InfoID: infoID,
		Uid:    principal,
	})

	if err != nil {
		log.Error(err)
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	type infoResultTag struct {
		Type string `json:"type"`
		Tag  string `json:"tag"`
	}

	type infoResult struct {
		Id              string            `json:"id"`
		UID             string            `json:"uid"`
		Author          string            `json:"author"`
		Title           string            `json:"title"`
		Summary         string            `json:"summary"`
		Country         string            `json:"country"`
		ContentTime     time.Time         `json:"contentTime"`
		CoverResources  map[string]string `json:"coverResources"`
		PublishTime     time.Time         `json:"publishTime"`
		LastReviewTime  time.Time         `json:"lastReviewTime"`
		Valid           bool              `json:"valid"`
		WatchCount      uint64            `json:"watchCount"`
		Tags            []*infoResultTag  `json:"tags"`
		ThumbUp         uint64            `json:"thumbUps"`
		IsThumbUp       bool              `json:"isThumbUp"`
		ThumbUpList     []string          `json:"thumbUpList"`
		ThumbDown       uint64            `json:"thumbDowns"`
		IsThumbDown     bool              `json:"isThumbDown"`
		ThumbDownList   []string          `json:"thumbDownList"`
		Favorites       uint64            `json:"favorites"`
		IsFavorite      bool              `json:"isFavorite"`
		FavoriteList    []string          `json:"favoriteList"`
		LastModifyTime  time.Time         `json:"lastModifyTime"`
		CanReview       bool              `json:"canReview"`
		Archived        bool              `json:"archived"`
		LatestSegmentID string            `json:"latestSegmentID"`
		SegmentCount    uint64            `json:"segmentCount"`
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

	tags := make([]*infoResultTag, 0, len(info.Tags))
	for _, v := range info.Tags {
		tags = append(tags, &infoResultTag{
			Type: v.Type,
			Tag:  v.Tag,
		})
	}

	resp := &infoResult{
		Id:              info.InfoID,
		UID:             info.Uid,
		Title:           info.Title,
		Author:          info.Author,
		Summary:         info.Summary,
		Country:         info.Country,
		ContentTime:     contentTime,
		CoverResources:  info.CoverResources,
		PublishTime:     publishTime,
		LastReviewTime:  lastReviewTime,
		Valid:           info.Valid,
		WatchCount:      info.WatchCount,
		Tags:            tags,
		ThumbUp:         info.ThumbUps,
		IsThumbUp:       info.IsThumbUp,
		ThumbUpList:     info.ThumbUpList,
		ThumbDown:       info.ThumbDowns,
		IsThumbDown:     info.IsThumbDown,
		ThumbDownList:   info.ThumbDownList,
		Favorites:       info.Favorites,
		IsFavorite:      info.IsFavorite,
		FavoriteList:    info.FavoriteList,
		LastModifyTime:  lastModifyTime,
		CanReview:       info.CanReview,
		Archived:        info.Archived,
		LatestSegmentID: info.LatestSegmentID,
		SegmentCount:    info.SegmentCount,
	}
	ctx.JSON(http.StatusOK, resp)
}

func (h *Content) UpdateInfo(ctx *gin.Context) {
	authPayload, err := h.middleware.ExtractToken(ctx)
	if err == gin_jwt.ErrContextNotHaveToken {
		ctx.Error(err)
		return
	}
	principal := authPayload["sub"].(string)

	contentClient, ok := clients.ContentFromContext(ctx)
	if !ok {
		ctx.AbortWithError(http.StatusInternalServerError, ErrClientNotFound)
		return
	}

	type updateInfoTag struct {
		Type string `json:"type"`
		Tag  string `json:"tag"`
	}

	type updateInfoReq struct {
		Title          string            `json:"title"`
		Country        string            `json:"country"`
		Tags           []*updateInfoTag  `json:"tags"`
		CanReview      bool              `json:"canReview"`
		CoverResources map[string]string `json:"coverResources"`
	}
	var req updateInfoReq
	ctx.Bind(&req)

	tags := make([]*content.TagAndType, 0, len(req.Tags))
	for _, v := range req.Tags {
		tags = append(tags, &content.TagAndType{
			Type: v.Type,
			Tag:  v.Tag,
		})
	}

	infoID := ctx.Param("id")
	_, err = contentClient.EditInfo(ctx, &content.EditInfoReq{
		Uid:            principal,
		InfoID:         infoID,
		Title:          req.Title,
		Country:        req.Country,
		Tags:           tags,
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

func (h *Content) DeleteInfo(ctx *gin.Context) {
	authPayload, err := h.middleware.ExtractToken(ctx)
	if err == gin_jwt.ErrContextNotHaveToken {
		ctx.Error(err)
		return
	}
	principal := authPayload["sub"].(string)

	contentClient, ok := clients.ContentFromContext(ctx)
	if !ok {
		ctx.AbortWithError(http.StatusInternalServerError, ErrClientNotFound)
		return
	}

	infoID := ctx.Param("id")

	_, err = contentClient.DeleteInfo(ctx, &content.InfoIDAndUIDReq{
		Uid:    principal,
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
	authPayload, err := h.middleware.ExtractToken(ctx)
	if err == gin_jwt.ErrContextNotHaveToken {
		ctx.Error(err)
		return
	}
	principal := authPayload["sub"].(string)

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
	var resp *content.InfoIDsResp
	if strings.Contains(ctx.Request.RequestURI, "favorite") {
		resp, err = contentClient.GetUserFavorite(ctx, &content.UIDPageReq{
			Page:  uint32(page),
			Size:  uint32(size),
			Sorts: sorts,
			Uid:   principal,
		})
	} else if strings.Contains(ctx.Request.RequestURI, "thumbUp") {
		resp, err = contentClient.GetUserThumbUp(ctx, &content.UIDPageReq{
			Page:  uint32(page),
			Size:  uint32(size),
			Sorts: sorts,
			Uid:   principal,
		})
	} else {
		resp, err = contentClient.GetUserThumbDown(ctx, &content.UIDPageReq{
			Page:  uint32(page),
			Size:  uint32(size),
			Sorts: sorts,
			Uid:   principal,
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
	//principal := ""
	//authPayload, err := h.middleware.ExtractToken(ctx)
	//if err != gin_jwt.ErrContextNotHaveToken {
	//	principal = authPayload["sub"].(string)
	//}

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
	authPayload, err := h.middleware.ExtractToken(ctx)
	if err == gin_jwt.ErrContextNotHaveToken {
		ctx.Error(err)
		return
	}
	principal := authPayload["sub"].(string)

	contentClient, ok := clients.ContentFromContext(ctx)
	if !ok {
		ctx.AbortWithError(http.StatusInternalServerError, ErrClientNotFound)
		return
	}

	infoID := ctx.Param("id")
	if strings.Contains(ctx.Request.RequestURI, "favorite") {
		_, err = contentClient.Favorite(ctx, &content.InfoIDAndUIDReq{
			InfoID: infoID,
			Uid:    principal,
		})
	} else if strings.Contains(ctx.Request.RequestURI, "thumbUp") {
		_, err = contentClient.ThumbUp(ctx, &content.InfoIDAndUIDReq{
			InfoID: infoID,
			Uid:    principal,
		})
	} else {
		_, err = contentClient.ThumbDown(ctx, &content.InfoIDAndUIDReq{
			InfoID: infoID,
			Uid:    principal,
		})
	}

	if err != nil {
		log.Error(err)
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	ctx.Status(http.StatusOK)
}

func (h *Content) DeleteFavThumb(ctx *gin.Context) {
	authPayload, err := h.middleware.ExtractToken(ctx)
	if err == gin_jwt.ErrContextNotHaveToken {
		ctx.Error(err)
		return
	}
	principal := authPayload["sub"].(string)

	contentClient, ok := clients.ContentFromContext(ctx)
	if !ok {
		ctx.AbortWithError(http.StatusInternalServerError, ErrClientNotFound)
		return
	}

	infoID := ctx.Param("id")
	if strings.Contains(ctx.Request.RequestURI, "favorite") {
		_, err = contentClient.DeleteFavorite(ctx, &content.InfoIDAndUIDReq{
			InfoID: infoID,
			Uid:    principal,
		})
	} else if strings.Contains(ctx.Request.RequestURI, "thumbUp") {
		_, err = contentClient.DeleteThumbUp(ctx, &content.InfoIDAndUIDReq{
			InfoID: infoID,
			Uid:    principal,
		})
	} else {
		_, err = contentClient.DeleteThumbDown(ctx, &content.InfoIDAndUIDReq{
			InfoID: infoID,
			Uid:    principal,
		})
	}

	if err != nil {
		log.Error(err)
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	ctx.Status(http.StatusOK)
}

func (h *Content) GetAllSegments(ctx *gin.Context) {
	principal := ""
	authPayload, err := h.middleware.ExtractToken(ctx)
	if err != gin_jwt.ErrContextNotHaveToken {
		principal = authPayload["sub"].(string)
	}

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
	if ctx.Query("sorts") != "" {
		sorts, err = buildSort(ctx.Query("sorts"))
		if err != nil {
			ctx.AbortWithError(http.StatusInternalServerError, err)
			return
		}
	}

	var labels []string
	if ctx.Query("labels") != "" {
		labels = strings.Split(ctx.Query("labels"), ",")
	}

	infoID := ctx.Param("id")
	resp, err := contentClient.GetSegments(ctx, &content.GetSegmentsReq{
		Uid:    principal,
		InfoID: infoID,
		Labels: labels,
		Page:   uint32(page),
		Size:   uint32(size),
		Sorts:  sorts,
	})

	if err != nil {
		log.Error(err)
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	type segmentResult struct {
		ID      string            `json:"id"`
		InfoID  string            `json:"infoID"`
		No      int64             `json:"no"`
		Title   string            `json:"title"`
		Labels  []string          `json:"labels"`
		Content map[string]string `json:"content"`
	}
	results := make([]*segmentResult, 0, len(resp.Segments))
	for _, seg := range resp.Segments {
		results = append(results, &segmentResult{
			ID:      seg.Id,
			InfoID:  seg.InfoID,
			No:      seg.No,
			Title:   seg.Title,
			Labels:  seg.Labels,
			Content: seg.Content,
		})
	}

	ctx.JSON(http.StatusOK, gin.H{
		"total_count": resp.TotalCount,
		"items":       results,
	})
}

func (h *Content) GetSegmentDetail(ctx *gin.Context) {
	principal := ""
	authPayload, err := h.middleware.ExtractToken(ctx)
	if err != gin_jwt.ErrContextNotHaveToken {
		principal = authPayload["sub"].(string)
	}

	contentClient, ok := clients.ContentFromContext(ctx)
	if !ok {
		ctx.AbortWithError(http.StatusInternalServerError, ErrClientNotFound)
		return
	}

	infoID := ctx.Param("id")
	segID := ctx.Param("segID")

	seg, err := contentClient.GetSegment(ctx, &content.InfoIDAndUIDAndSegIDReq{
		Uid:    principal,
		InfoID: infoID,
		SegID:  segID,
	})

	if err != nil {
		log.Error(err)
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	type segmentResult struct {
		ID      string            `json:"id"`
		InfoID  string            `json:"infoID"`
		No      int64             `json:"no"`
		Title   string            `json:"title"`
		Labels  []string          `json:"labels"`
		Content map[string]string `json:"content"`
	}

	resp := &segmentResult{
		ID:      seg.Id,
		InfoID:  seg.InfoID,
		No:      seg.No,
		Title:   seg.Title,
		Labels:  seg.Labels,
		Content: seg.Content,
	}
	ctx.JSON(http.StatusOK, resp)
}

func (h *Content) PublishSegment(ctx *gin.Context) {
	authPayload, err := h.middleware.ExtractToken(ctx)
	if err == gin_jwt.ErrContextNotHaveToken {
		ctx.Error(err)
		return
	}
	principal := authPayload["sub"].(string)

	contentClient, ok := clients.ContentFromContext(ctx)
	if !ok {
		ctx.AbortWithError(http.StatusInternalServerError, ErrClientNotFound)
		return
	}

	type publishSegmentReq struct {
		InfoID  string            `json:"infoID"`
		No      int64             `json:"no"`
		Labels  []string          `json:"labels"`
		Content map[string]string `json:"content"`
	}
	var req publishSegmentReq
	ctx.Bind(&req)

	_, err = contentClient.PublishSegment(ctx, &content.PublishSegmentReq{
		Uid:     principal,
		InfoID:  req.InfoID,
		No:      req.No,
		Labels:  req.Labels,
		Content: req.Content,
	})

	if err != nil {
		log.Error(err)
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	ctx.Status(http.StatusOK)
}

func (h *Content) UpdateSegment(ctx *gin.Context) {
	authPayload, err := h.middleware.ExtractToken(ctx)
	if err == gin_jwt.ErrContextNotHaveToken {
		ctx.Error(err)
		return
	}
	principal := authPayload["sub"].(string)

	contentClient, ok := clients.ContentFromContext(ctx)
	if !ok {
		ctx.AbortWithError(http.StatusInternalServerError, ErrClientNotFound)
		return
	}

	type updateSegmentReq struct {
		No      int64             `json:"no"`
		Labels  []string          `json:"labels"`
		Content map[string]string `json:"content"`
	}
	var req updateSegmentReq
	ctx.Bind(&req)

	infoID := ctx.Param("id")
	segID := ctx.Param("segID")

	_, err = contentClient.EditSegment(ctx, &content.EditSegmentReq{
		Uid:     principal,
		InfoID:  infoID,
		SegID:   segID,
		No:      req.No,
		Labels:  req.Labels,
		Content: req.Content,
	})

	if err != nil {
		log.Error(err)
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	ctx.Status(http.StatusOK)
}

func (h *Content) DeleteSegment(ctx *gin.Context) {
	authPayload, err := h.middleware.ExtractToken(ctx)
	if err == gin_jwt.ErrContextNotHaveToken {
		ctx.Error(err)
		return
	}
	principal := authPayload["sub"].(string)

	contentClient, ok := clients.ContentFromContext(ctx)
	if !ok {
		ctx.AbortWithError(http.StatusInternalServerError, ErrClientNotFound)
		return
	}

	infoID := ctx.Param("id")
	segID := ctx.Param("segID")

	_, err = contentClient.DeleteSegment(ctx, &content.InfoIDAndUIDAndSegIDReq{
		Uid:    principal,
		InfoID: infoID,
		SegID:  segID,
	})

	if err != nil {
		log.Error(err)
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	ctx.Status(http.StatusOK)
}
