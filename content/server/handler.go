package server

import (
	"github.com/golang/protobuf/ptypes/empty"
	"github.com/golang/protobuf/ptypes/timestamp"
	"github.com/mongodb/mongo-go-driver/bson/objectid"
	"github.com/mongodb/mongo-go-driver/mongo"
	log "github.com/sirupsen/logrus"
	"github.com/zhsyourai/teddy-backend/common/proto/content"
	"github.com/zhsyourai/teddy-backend/content/models"
	"github.com/zhsyourai/teddy-backend/content/repositories"
	"golang.org/x/net/context"
	"time"
)

func NewContentServer(client *mongo.Client) (content.ContentServer, error) {
	// New Repository
	infoRepo, err := repositories.NewInfoRepository(client)
	if err != nil {
		return nil, err
	}
	favoriteRepo, err := repositories.NewFavoriteRepository(client)
	if err != nil {
		return nil, err
	}
	tagRepo, err := repositories.NewTagRepository(client)
	if err != nil {
		return nil, err
	}
	thumbUpRepo, err := repositories.NewThumbUpRepository(client)
	if err != nil {
		return nil, err
	}
	thumbDownRepo, err := repositories.NewThumbDownRepository(client)
	if err != nil {
		return nil, err
	}
	instance := &contentHandler{
		client:        client,
		infoRepo:      infoRepo,
		tagRepo:       tagRepo,
		favoriteRepo:  favoriteRepo,
		thumbDownRepo: thumbDownRepo,
		thumbUpRepo:   thumbUpRepo,
	}
	return instance, nil
}

type contentHandler struct {
	client        *mongo.Client
	infoRepo      repositories.InfoRepository
	tagRepo       repositories.TagRepository
	favoriteRepo  repositories.BehaviorRepository
	thumbUpRepo   repositories.BehaviorRepository
	thumbDownRepo repositories.BehaviorRepository
}

func (h *contentHandler) GetTags(ctx context.Context, req *content.GetTagReq) (*content.GetTagsResp, error) {
	var resp content.GetTagsResp
	if err := validateGetTagsReq(req); err != nil {
		return nil, err
	}

	err := h.client.UseSession(ctx, func(sessionContext mongo.SessionContext) error {
		var err error
		tags, err := h.tagRepo.FindAll(sessionContext, req.Type, req.Page, req.Size, req.Sorts)

		if err != nil {
			return err
		}

		result := make([]*content.Tag, 0, len(tags))
		for _, tag := range tags {
			pbTag := &content.Tag{}
			copyFromTagToPBTag(tag, pbTag)
			result = append(result, pbTag)
		}
		resp.Tags = result
		return nil
	})

	if err != nil {
		return nil, err
	}
	return &resp, nil
}

func (h *contentHandler) PublishInfo(ctx context.Context, req *content.PublishInfoReq) (*empty.Empty, error) {
	var resp empty.Empty
	if err := validatePublishInfoReq(req); err != nil {
		return nil, err
	}

	timeoutCtx, cancel := context.WithTimeout(ctx, 20*time.Second)
	defer cancel()

	err := h.client.UseSession(timeoutCtx, func(sessionContext mongo.SessionContext) error {
		var err error
		now := time.Now()
		err = sessionContext.StartTransaction()
		if err != nil {
			return err
		}

		defer func() {
			if err == nil {
				err = sessionContext.CommitTransaction(sessionContext)
			} else {
				err = sessionContext.AbortTransaction(context.Background())
			}
		}()

		tags := make([]*models.TypeAndTag, 0, len(req.Tags))
		for _, v := range req.Tags {
			result, err := h.tagRepo.FindByTypeAndTag(sessionContext, v.Type, v.Tag)
			if err != nil {
				log.Errorf("tag can't find error %v", err)
				return err
			}
			tags = append(tags, &models.TypeAndTag{
				Type: result.Type,
				Tag:  result.Tag,
			})
			err = h.tagRepo.IncUsage(sessionContext, result.ID, 1)
			if err != nil {
				log.Errorf("inc tag usage error %v", err)
				return err
			}
		}

		info := models.Info{
			Id:             objectid.New(),
			UID:            req.Uid,
			Title:          req.Title,
			Content:        req.Content,
			CoverResources: req.CoverResources,
			PublishTime:    now,
			LastReviewTime: time.Now(),
			Valid:          true,
			WatchCount:     0,
			Tags:           tags,
			ThumbUp:        0,
			ThumbDown:      0,
			Favorites:      0,
			LastModifyTime: now,
			CanReview:      req.CanReview,
			Archived:       false,
		}

		err = h.infoRepo.Insert(sessionContext, &info)
		if err != nil {
			return err
		}
		return nil
	})

	if err != nil {
		return nil, err
	}
	return &resp, nil
}

func difference(a, b []objectid.ObjectID) ([]objectid.ObjectID, []objectid.ObjectID) {
	ma := map[objectid.ObjectID]bool{}
	for _, x := range a {
		ma[x] = true
	}
	mb := map[objectid.ObjectID]bool{}
	for _, x := range b {
		mb[x] = true
	}
	ab := make([]objectid.ObjectID, 0, 20)
	for _, x := range a {
		if _, ok := mb[x]; !ok {
			ab = append(ab, x)
		}
	}
	ba := make([]objectid.ObjectID, 0, 20)
	for _, x := range a {
		if _, ok := mb[x]; !ok {
			ab = append(ab, x)
		}
	}
	return ab, ba
}

func (h *contentHandler) EditInfo(ctx context.Context, req *content.EditInfoReq) (*empty.Empty, error) {
	var resp empty.Empty
	if err := validateEditInfoReq(req); err != nil {
		return nil, err
	}

	infoID, err := objectid.FromHex(req.InfoID)
	if err != nil {
		return nil, err
	}

	timeoutCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	err = h.client.UseSession(timeoutCtx, func(sessionContext mongo.SessionContext) error {
		var err error
		now := time.Now()
		err = sessionContext.StartTransaction()
		if err != nil {
			return err
		}

		defer func() {
			if err == nil {
				err = sessionContext.CommitTransaction(sessionContext)
			} else {
				err = sessionContext.AbortTransaction(context.Background())
			}
		}()

		curInfo, err := h.infoRepo.FindOne(sessionContext, infoID)
		if err != nil {
			return err
		}

		curTagIDs := make([]objectid.ObjectID, 0, len(curInfo.Tags))
		for _, v := range curInfo.Tags {
			result, err := h.tagRepo.FindByTypeAndTag(sessionContext, v.Type, v.Tag)
			if err != nil {
				log.Errorf("tag can't find error %v", err)
				return err
			}
			curTagIDs = append(curTagIDs, result.ID)
		}

		tagIDs := make([]objectid.ObjectID, 0, len(req.Tags))
		tags := make([]*models.TypeAndTag, 0, len(req.Tags))
		for _, v := range req.Tags {
			result, err := h.tagRepo.FindByTypeAndTag(sessionContext, v.Type, v.Tag)
			if err != nil {
				log.Errorf("tag can't find error %v", err)
				return err
			}
			tagIDs = append(tagIDs, result.ID)
			tags = append(tags, &models.TypeAndTag{
				Tag:  result.Tag,
				Type: result.Type,
			})
		}

		sub, inc := difference(curTagIDs, tagIDs)

		for _, v := range inc {
			err = h.tagRepo.IncUsage(sessionContext, v, 1)
			if err != nil {
				log.Errorf("inc tag usage error %v", err)
				return err
			}
		}

		for _, v := range sub {
			err = h.tagRepo.IncUsage(sessionContext, v, -1)
			if err != nil {
				log.Errorf("sub tag usage error %v", err)
				return err
			}
		}

		err = h.infoRepo.Update(sessionContext, infoID, map[string]interface{}{
			"uid":            req.Uid,
			"title":          req.Title,
			"author":         req.Author,
			"summary":        req.Summary,
			"content":        req.Content,
			"coverResources": req.CoverResources,
			"canReview":      req.CanReview,
			"tags":           tags,
			"lastModifyTime": now,
		})
		if err != nil {
			return err
		}
		return nil
	})

	if err != nil {
		return nil, err
	}
	return &resp, nil
}

func (h *contentHandler) fillInfo(sessionContext mongo.SessionContext, uid string, info *models.Info) (*content.Info, error) {
	isThumbUp, err := h.thumbUpRepo.IsExists(sessionContext, uid, info.Id)
	if err != nil {
		return nil, err
	}

	isThumbDown, err := h.thumbDownRepo.IsExists(sessionContext, uid, info.Id)
	if err != nil {
		return nil, err
	}

	isFavorite, err := h.favoriteRepo.IsExists(sessionContext, uid, info.Id)
	if err != nil {
		return nil, err
	}

	tmpList, err := h.thumbUpRepo.FindUserByInfo(sessionContext, info.Id, 0, 10, []*content.Sort{
		&content.Sort{
			Name: "time",
			Asc:  false,
		},
	})
	if err != nil {
		return nil, err
	}
	thumbUpList := make([]string, 0, len(tmpList))
	for _, v := range tmpList {
		thumbUpList = append(thumbUpList, v.UID)
	}

	tmpList, err = h.thumbDownRepo.FindUserByInfo(sessionContext, info.Id, 0, 10, []*content.Sort{
		&content.Sort{
			Name: "time",
			Asc:  false,
		},
	})
	if err != nil {
		return nil, err
	}
	thumbDownList := make([]string, 0, len(tmpList))
	for _, v := range tmpList {
		thumbDownList = append(thumbDownList, v.UID)
	}

	tmpList, err = h.favoriteRepo.FindUserByInfo(sessionContext, info.Id, 0, 10, []*content.Sort{
		&content.Sort{
			Name: "time",
			Asc:  false,
		},
	})
	if err != nil {
		return nil, err
	}
	favoriteList := make([]string, 0, len(tmpList))
	for _, v := range tmpList {
		favoriteList = append(favoriteList, v.UID)
	}

	tagList := make([]*content.TagAndType, 0, len(info.Tags))
	for _, v := range info.Tags {
		tagList = append(tagList, &content.TagAndType{
			Tag:  v.Tag,
			Type: v.Type,
		})
	}

	resp := &content.Info{
		InfoID:  info.Id.Hex(),
		Uid:     info.UID,
		Title:   info.Title,
		Content: info.Content,
		ContentTime: &timestamp.Timestamp{
			Seconds: info.ContentTime.Unix(),
			Nanos:   int32(info.ContentTime.Nanosecond()),
		},
		CoverResources: info.CoverResources,
		PublishTime: &timestamp.Timestamp{
			Seconds: info.PublishTime.Unix(),
			Nanos:   int32(info.PublishTime.Nanosecond()),
		},
		LastReviewTime: &timestamp.Timestamp{
			Seconds: info.LastReviewTime.Unix(),
			Nanos:   int32(info.LastReviewTime.Nanosecond()),
		},
		Valid:         info.Valid,
		WatchCount:    info.WatchCount,
		Tags:          tagList,
		ThumbUp:       info.ThumbUp,
		IsThumbUp:     isThumbUp,
		ThumbUpList:   thumbUpList,
		ThumbDown:     info.ThumbDown,
		IsThumbDown:   isThumbDown,
		ThumbDownList: thumbDownList,
		Favorites:     info.Favorites,
		IsFavorite:    isFavorite,
		FavoriteList:  favoriteList,
		LastModifyTime: &timestamp.Timestamp{
			Seconds: info.LastModifyTime.Unix(),
			Nanos:   int32(info.LastModifyTime.Nanosecond()),
		},
		CanReview: info.CanReview,
	}
	return resp, nil
}

func (h *contentHandler) GetInfo(ctx context.Context, req *content.GetInfoReq) (*content.Info, error) {
	var resp *content.Info
	var err error
	if err = validateGetInfoReq(req); err != nil {
		return nil, err
	}

	infoID, err := objectid.FromHex(req.InfoID)
	if err != nil {
		return nil, err
	}

	err = h.client.UseSession(ctx, func(sessionContext mongo.SessionContext) error {
		info, err := h.infoRepo.FindOne(sessionContext, infoID)
		if err != nil {
			return err
		}

		resp, err = h.fillInfo(sessionContext, req.Uid, info)
		if err != nil {
			return err
		}
		return nil
	})

	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (h *contentHandler) GetInfos(ctx context.Context, req *content.GetInfosReq) (*content.GetInfosResp, error) {
	var resp content.GetInfosResp
	if err := validateGetInfosReq(req); err != nil {
		return nil, err
	}

	timeoutCtx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	err := h.client.UseSession(timeoutCtx, func(sessionContext mongo.SessionContext) error {
		var err error
		tags := make([]*models.TypeAndTag, 0, len(req.Tags))
		for _, v := range req.Tags {
			tags = append(tags, &models.TypeAndTag{
				Type: v.Type,
				Tag:  v.Tag,
			})
		}

		var infos []*models.Info
		if req.Title == "" {
			infos, err = h.infoRepo.FindAll(sessionContext, req.Uid, tags, req.Page, req.Size, req.Sorts)
		} else {
			//TODO: search use elasticsearch
		}
		if err != nil {
			return err
		}

		pInfos := make([]*content.Info, 0, len(infos))
		for _, info := range infos {
			pInfo, err := h.fillInfo(sessionContext, req.Uid, info)
			if err != nil {
				return err
			}
			pInfos = append(pInfos, pInfo)
		}
		resp.Infos = pInfos
		return nil
	})

	if err != nil {
		return nil, err
	}
	return &resp, nil
}

func (h *contentHandler) DeleteInfo(ctx context.Context, req *content.InfoIDReq) (*empty.Empty, error) {
	var resp empty.Empty
	if err := validateInfoIDReq(req); err != nil {
		return nil, err
	}

	err := h.client.UseSession(ctx, func(sessionContext mongo.SessionContext) error {
		var err error
		infoID, err := objectid.FromHex(req.InfoID)
		if err != nil {
			return err
		}

		err = h.infoRepo.Delete(sessionContext, infoID)
		if err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return nil, err
	}
	return &resp, nil
}

func (h *contentHandler) WatchInfo(ctx context.Context, req *content.InfoIDReq) (*empty.Empty, error) {
	var resp empty.Empty
	if err := validateInfoIDReq(req); err != nil {
		return nil, err
	}

	err := h.client.UseSession(ctx, func(sessionContext mongo.SessionContext) error {
		var err error
		infoID, err := objectid.FromHex(req.InfoID)
		if err != nil {
			return err
		}

		err = h.infoRepo.IncWatchCount(sessionContext, infoID, 1)
		if err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return nil, err
	}
	return &resp, nil
}

func (h *contentHandler) _behaviorInsert(ctx context.Context,
	checkRepo, repo repositories.BehaviorRepository, req *content.InfoIDAndUIDReq) (*empty.Empty, error) {
	var resp empty.Empty

	if err := validateInfoIDAndUIDReq(req); err != nil {
		return nil, err
	}

	infoID, err := objectid.FromHex(req.InfoID)
	if err != nil {
		return nil, err
	}

	err = h.client.UseSession(ctx, func(sessionContext mongo.SessionContext) error {
		var err error
		isExist, err := checkRepo.IsExists(sessionContext, req.Uid, infoID)
		if err != nil {
			return err
		}

		if isExist {
			return ErrBehaviorExists
		}

		err = repo.Insert(sessionContext, req.Uid, &models.BehaviorInfoItem{
			InfoId: infoID,
			Time:   time.Now(),
		})
		if err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return &resp, nil
}

func (h *contentHandler) _behaviorFindInfoByUser(ctx context.Context,
	repo repositories.BehaviorRepository, req *content.UIDPageReq) (*content.InfoIDsResp, error) {
	var resp content.InfoIDsResp
	if err := validateUIDPageReq(req); err != nil {
		return nil, err
	}

	err := h.client.UseSession(ctx, func(sessionContext mongo.SessionContext) error {
		items, err := repo.FindInfoByUser(sessionContext, req.Uid, req.Page, req.Size, req.Sorts)

		if err != nil {
			return err
		}

		results := make([]*content.InfoIDWithTime, 0, len(items))
		for _, item := range items {
			results = append(results, &content.InfoIDWithTime{
				InfoId: item.InfoId.Hex(),
				Time: &timestamp.Timestamp{
					Seconds: item.Time.Unix(),
					Nanos:   int32(item.Time.Nanosecond()),
				},
			})
		}
		resp.Ids = results
		return nil
	})

	if err != nil {
		return nil, err
	}

	return &resp, nil
}

func (h *contentHandler) _behaviorFindUserByInfo(ctx context.Context,
	repo repositories.BehaviorRepository, req *content.InfoIDPageReq) (*content.UserIDsResp, error) {
	var resp content.UserIDsResp
	if err := validateInfoIDPageReq(req); err != nil {
		return nil, err
	}

	infoID, err := objectid.FromHex(req.InfoID)
	if err != nil {
		return nil, err
	}

	err = h.client.UseSession(ctx, func(sessionContext mongo.SessionContext) error {
		items, err := repo.FindUserByInfo(sessionContext, infoID, req.Page, req.Size, req.Sorts)

		if err != nil {
			return err
		}

		results := make([]*content.UIDWithTime, 0, len(items))
		for _, item := range items {
			results = append(results, &content.UIDWithTime{
				Uid: item.UID,
				Time: &timestamp.Timestamp{
					Seconds: item.Time.Unix(),
					Nanos:   int32(item.Time.Nanosecond()),
				},
			})
		}
		resp.Items = results
		return nil
	})

	if err != nil {
		return nil, err
	}

	return &resp, nil
}

func (h *contentHandler) _behaviorDelete(ctx context.Context,
	repo repositories.BehaviorRepository, req *content.InfoIDAndUIDReq) (*empty.Empty, error) {
	var resp empty.Empty
	if err := validateInfoIDAndUIDReq(req); err != nil {
		return nil, err
	}

	infoID, err := objectid.FromHex(req.InfoID)
	if err != nil {
		return nil, err
	}

	err = h.client.UseSession(ctx, func(sessionContext mongo.SessionContext) error {
		err = repo.Delete(sessionContext, req.Uid, infoID)
		if err != nil {
			return err
		}
		return nil
	})

	if err != nil {
		return nil, err
	}
	return &resp, nil
}

func (h *contentHandler) ThumbUp(ctx context.Context, req *content.InfoIDAndUIDReq) (*empty.Empty, error) {
	return h._behaviorInsert(ctx, h.thumbDownRepo, h.thumbUpRepo, req)
}

func (h *contentHandler) ThumbDown(ctx context.Context, req *content.InfoIDAndUIDReq) (*empty.Empty, error) {
	return h._behaviorInsert(ctx, h.thumbUpRepo, h.thumbDownRepo, req)
}

func (h *contentHandler) GetUserThumbUp(ctx context.Context, req *content.UIDPageReq) (*content.InfoIDsResp, error) {
	return h._behaviorFindInfoByUser(ctx, h.thumbUpRepo, req)
}

func (h *contentHandler) GetUserThumbDown(ctx context.Context, req *content.UIDPageReq) (*content.InfoIDsResp, error) {
	return h._behaviorFindInfoByUser(ctx, h.thumbDownRepo, req)
}

func (h *contentHandler) GetInfoThumbUp(ctx context.Context, req *content.InfoIDPageReq) (*content.UserIDsResp, error) {
	return h._behaviorFindUserByInfo(ctx, h.thumbUpRepo, req)
}

func (h *contentHandler) GetInfoThumbDown(ctx context.Context, req *content.InfoIDPageReq) (*content.UserIDsResp, error) {
	return h._behaviorFindUserByInfo(ctx, h.thumbDownRepo, req)
}

func (h *contentHandler) DeleteThumbUp(ctx context.Context, req *content.InfoIDAndUIDReq) (*empty.Empty, error) {
	return h._behaviorDelete(ctx, h.thumbUpRepo, req)
}

func (h *contentHandler) DeleteThumbDown(ctx context.Context, req *content.InfoIDAndUIDReq) (*empty.Empty, error) {
	return h._behaviorDelete(ctx, h.thumbDownRepo, req)
}

func (h *contentHandler) Favorite(ctx context.Context, req *content.InfoIDAndUIDReq) (*empty.Empty, error) {
	return h._behaviorInsert(ctx, h.favoriteRepo, h.favoriteRepo, req)
}

func (h *contentHandler) GetUserFavorite(ctx context.Context, req *content.UIDPageReq) (*content.InfoIDsResp, error) {
	return h._behaviorFindInfoByUser(ctx, h.favoriteRepo, req)
}

func (h *contentHandler) GetInfoFavorite(ctx context.Context, req *content.InfoIDPageReq) (*content.UserIDsResp, error) {
	return h._behaviorFindUserByInfo(ctx, h.favoriteRepo, req)
}

func (h *contentHandler) DeleteFavorite(ctx context.Context, req *content.InfoIDAndUIDReq) (*empty.Empty, error) {
	return h._behaviorDelete(ctx, h.favoriteRepo, req)
}
