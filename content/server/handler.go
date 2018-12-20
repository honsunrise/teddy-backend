package server

import (
	"github.com/golang/protobuf/ptypes/empty"
	"github.com/golang/protobuf/ptypes/timestamp"
	"github.com/mongodb/mongo-go-driver/bson/objectid"
	log "github.com/sirupsen/logrus"
	"github.com/zhsyourai/teddy-backend/common/proto"
	"github.com/zhsyourai/teddy-backend/common/types"
	"github.com/zhsyourai/teddy-backend/content/converter"
	"github.com/zhsyourai/teddy-backend/content/models"
	"github.com/zhsyourai/teddy-backend/content/repositories"
	"golang.org/x/net/context"
	"time"
)

func NewContentServer(infoRepo repositories.InfoRepository,
	tagRepo repositories.TagRepository, favoriteRepo repositories.BehaviorRepository,
	thumbUpRepo repositories.BehaviorRepository, thumbDownRepo repositories.BehaviorRepository) (proto.ContentServer, error) {
	instance := &contentHandler{
		infoRepo:      infoRepo,
		tagRepo:       tagRepo,
		favoriteRepo:  favoriteRepo,
		thumbDownRepo: thumbDownRepo,
		thumbUpRepo:   thumbUpRepo,
	}
	return instance, nil
}

type contentHandler struct {
	infoRepo      repositories.InfoRepository
	tagRepo       repositories.TagRepository
	favoriteRepo  repositories.BehaviorRepository
	thumbUpRepo   repositories.BehaviorRepository
	thumbDownRepo repositories.BehaviorRepository
}

func (h *contentHandler) GetTags(ctx context.Context, req *proto.GetTagReq) (*proto.GetTagsResp, error) {
	var resp proto.GetTagsResp
	if err := validateGetTagsReq(req); err != nil {
		return nil, err
	}

	tags, err := h.tagRepo.FindAll(req.Page, req.Size, []types.Sort{
		{"usage", types.DESC},
	})

	if err != nil {
		return nil, err
	}

	result := make([]*proto.Tag, 0, len(tags))
	for _, tag := range tags {
		pbTag := &proto.Tag{}
		converter.CopyFromTagToPBTag(tag, pbTag)
		result = append(result, pbTag)
	}

	resp.Tags = result
	return &resp, nil
}

func (h *contentHandler) PublishInfo(ctx context.Context, req *proto.PublishInfoReq) (*empty.Empty, error) {
	var resp empty.Empty
	if err := validatePublishInfoReq(req); err != nil {
		return nil, err
	}
	now := time.Now()

	for _, tag := range req.Tags {
		_, err := h.tagRepo.FindOne(tag)
		if err == nil {
			err = h.tagRepo.IncUsage(tag, 1)
			if err != nil {
				log.Errorf("inc tag use times error %v", err)
			}
		} else {
			err = h.tagRepo.Insert(&models.Tag{
				Tag:         tag,
				Usage:       1,
				CreateTime:  now,
				LastUseTime: now,
			})
			if err != nil {
				log.Errorf("insert tag error %v", err)
			}
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
		Tags:           req.Tags,
		ThumbUp:        0,
		ThumbDown:      0,
		Favorites:      0,
		LastModifyTime: now,
		CanReview:      req.CanReview,
	}

	err := h.infoRepo.Insert(&info)
	if err != nil {
		return nil, err
	}

	return &resp, nil
}

func (h *contentHandler) EditInfo(ctx context.Context, req *proto.EditInfoReq) (*empty.Empty, error) {
	var resp empty.Empty
	if err := validateEditInfoReq(req); err != nil {
		return nil, err
	}

	infoID, err := objectid.FromHex(req.InfoID)
	if err != nil {
		return nil, err
	}

	now := time.Now()

	for _, tag := range req.Tags {
		_, err := h.tagRepo.FindOne(tag)
		if err == nil {
			err = h.tagRepo.IncUsage(tag, 1)
			if err != nil {
				log.Errorf("inc tag use times error %v", err)
			}
		} else {
			err = h.tagRepo.Insert(&models.Tag{
				Tag:         tag,
				Usage:       1,
				CreateTime:  now,
				LastUseTime: now,
			})
			if err != nil {
				log.Errorf("insert tag error %v", err)
			}
		}
	}

	err = h.infoRepo.Update(infoID, map[string]interface{}{
		"uid":            req.Uid,
		"title":          req.Title,
		"content":        req.Content,
		"coverResources": req.CoverResources,
		"canReview":      req.CanReview,
	})
	if err != nil {
		return nil, err
	}

	return &resp, nil
}

func (h *contentHandler) fillInfo(uid string, info *models.Info) (*proto.Info, error) {
	isThumbUp, err := h.thumbUpRepo.IsExist(uid, info.Id)
	if err != nil {
		return nil, err
	}

	isThumbDown, err := h.thumbDownRepo.IsExist(uid, info.Id)
	if err != nil {
		return nil, err
	}

	isFavorite, err := h.favoriteRepo.IsExist(uid, info.Id)
	if err != nil {
		return nil, err
	}

	tmpList, err := h.thumbUpRepo.FindUserByInfo(info.Id, 0, 10, []*proto.Sort{
		&proto.Sort{
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

	tmpList, err = h.thumbDownRepo.FindUserByInfo(info.Id, 0, 10, []*proto.Sort{
		&proto.Sort{
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

	tmpList, err = h.favoriteRepo.FindUserByInfo(info.Id, 0, 10, []*proto.Sort{
		&proto.Sort{
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

	resp := &proto.Info{
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
		Tags:          info.Tags,
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

func (h *contentHandler) GetInfo(ctx context.Context, req *proto.GetInfoReq) (*proto.Info, error) {
	if err := validateGetInfoReq(req); err != nil {
		return nil, err
	}

	infoID, err := objectid.FromHex(req.InfoID)
	if err != nil {
		return nil, err
	}

	info, err := h.infoRepo.FindOne(infoID)
	if err != nil {
		return nil, err
	}

	resp, err := h.fillInfo(req.Uid, info)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (h *contentHandler) GetInfos(ctx context.Context, req *proto.GetInfosReq) (*proto.GetInfosResp, error) {
	var resp proto.GetInfosResp
	if err := validateGetInfosReq(req); err != nil {
		return nil, err
	}

	var infos []*models.Info
	var err error

	if req.Uid != "" && req.Title == "" {
		infos, err = h.infoRepo.FindByUser(req.Uid, req.Page, req.Size, req.Sorts)
	} else if req.Uid == "" && req.Title != "" {
		//TODO: search use elasticsearch
	} else if req.Uid != "" && req.Title != "" {
		//TODO: search use elasticsearch
	} else {
		infos, err = h.infoRepo.FindAll(req.Page, req.Size, req.Sorts)
	}
	if err != nil {
		return nil, err
	}

	pInfos := make([]*proto.Info, 0, len(infos))
	for _, info := range infos {
		pInfo, err := h.fillInfo(req.Uid, info)
		if err != nil {
			return nil, err
		}
		pInfos = append(pInfos, pInfo)
	}
	resp.Infos = pInfos
	return &resp, nil
}

func (h *contentHandler) DeleteInfo(ctx context.Context, req *proto.InfoIDReq) (*empty.Empty, error) {
	var resp empty.Empty
	if err := validateInfoIDReq(req); err != nil {
		return nil, err
	}

	infoID, err := objectid.FromHex(req.InfoID)
	if err != nil {
		return nil, err
	}

	err = h.infoRepo.Delete(infoID)
	if err != nil {
		return nil, err
	}

	return &resp, nil
}

func (h *contentHandler) WatchInfo(ctx context.Context, req *proto.InfoIDReq) (*empty.Empty, error) {
	var resp empty.Empty
	if err := validateInfoIDReq(req); err != nil {
		return nil, err
	}

	infoID, err := objectid.FromHex(req.InfoID)
	if err != nil {
		return nil, err
	}
	err = h.infoRepo.IncWatchCount(infoID, 1)
	if err != nil {
		return nil, err
	}

	return &resp, nil
}

func (h *contentHandler) _behaviorInsert(ctx context.Context,
	repo repositories.BehaviorRepository, req *proto.InfoIDAndUIDReq) (*empty.Empty, error) {
	var resp empty.Empty
	if err := validateInfoIDAndUIDReq(req); err != nil {
		return nil, err
	}

	infoID, err := objectid.FromHex(req.InfoID)
	if err != nil {
		return nil, err
	}

	err = repo.Insert(req.Uid, &models.BehaviorInfoItem{
		InfoId: infoID,
		Time:   time.Now(),
	})

	if err != nil {
		return nil, err
	}

	return &resp, nil
}

func (h *contentHandler) _behaviorFindInfoByUser(ctx context.Context,
	repo repositories.BehaviorRepository, req *proto.UIDPageReq) (*proto.InfoIDsResp, error) {
	var resp proto.InfoIDsResp
	if err := validateUIDPageReq(req); err != nil {
		return nil, err
	}

	items, err := repo.FindInfoByUser(req.Uid, req.Page, req.Size, req.Sorts)

	if err != nil {
		return nil, err
	}

	results := make([]*proto.InfoIDWithTime, 0, len(items))
	for _, item := range items {
		results = append(results, &proto.InfoIDWithTime{
			InfoId: item.InfoId.Hex(),
			Time: &timestamp.Timestamp{
				Seconds: item.Time.Unix(),
				Nanos:   int32(item.Time.Nanosecond()),
			},
		})
	}
	resp.Ids = results
	return &resp, nil
}

func (h *contentHandler) _behaviorFindUserByInfo(ctx context.Context,
	repo repositories.BehaviorRepository, req *proto.InfoIDPageReq) (*proto.UserIDsResp, error) {
	var resp proto.UserIDsResp
	if err := validateInfoIDPageReq(req); err != nil {
		return nil, err
	}

	infoID, err := objectid.FromHex(req.InfoID)
	if err != nil {
		return nil, err
	}

	items, err := repo.FindUserByInfo(infoID, req.Page, req.Size, req.Sorts)

	if err != nil {
		return nil, err
	}

	results := make([]*proto.UIDWithTime, 0, len(items))
	for _, item := range items {
		results = append(results, &proto.UIDWithTime{
			Uid: item.UID,
			Time: &timestamp.Timestamp{
				Seconds: item.Time.Unix(),
				Nanos:   int32(item.Time.Nanosecond()),
			},
		})
	}
	resp.Items = results
	return &resp, nil
}

func (h *contentHandler) _behaviorDelete(ctx context.Context,
	repo repositories.BehaviorRepository, req *proto.InfoIDAndUIDReq) (*empty.Empty, error) {
	var resp empty.Empty
	if err := validateInfoIDAndUIDReq(req); err != nil {
		return nil, err
	}

	infoID, err := objectid.FromHex(req.InfoID)
	if err != nil {
		return nil, err
	}

	err = repo.Delete(req.Uid, infoID)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

func (h *contentHandler) ThumbUp(ctx context.Context, req *proto.InfoIDAndUIDReq) (*empty.Empty, error) {
	return h._behaviorInsert(ctx, h.thumbUpRepo, req)
}

func (h *contentHandler) ThumbDown(ctx context.Context, req *proto.InfoIDAndUIDReq) (*empty.Empty, error) {
	return h._behaviorInsert(ctx, h.thumbDownRepo, req)
}

func (h *contentHandler) GetUserThumbUp(ctx context.Context, req *proto.UIDPageReq) (*proto.InfoIDsResp, error) {
	return h._behaviorFindInfoByUser(ctx, h.thumbUpRepo, req)
}

func (h *contentHandler) GetUserThumbDown(ctx context.Context, req *proto.UIDPageReq) (*proto.InfoIDsResp, error) {
	return h._behaviorFindInfoByUser(ctx, h.thumbDownRepo, req)
}

func (h *contentHandler) GetInfoThumbUp(ctx context.Context, req *proto.InfoIDPageReq) (*proto.UserIDsResp, error) {
	return h._behaviorFindUserByInfo(ctx, h.thumbUpRepo, req)
}

func (h *contentHandler) GetInfoThumbDown(ctx context.Context, req *proto.InfoIDPageReq) (*proto.UserIDsResp, error) {
	return h._behaviorFindUserByInfo(ctx, h.thumbDownRepo, req)
}

func (h *contentHandler) DeleteThumbUp(ctx context.Context, req *proto.InfoIDAndUIDReq) (*empty.Empty, error) {
	return h._behaviorDelete(ctx, h.thumbUpRepo, req)
}

func (h *contentHandler) DeleteThumbDown(ctx context.Context, req *proto.InfoIDAndUIDReq) (*empty.Empty, error) {
	return h._behaviorDelete(ctx, h.thumbDownRepo, req)
}

func (h *contentHandler) Favorite(ctx context.Context, req *proto.InfoIDAndUIDReq) (*empty.Empty, error) {
	return h._behaviorInsert(ctx, h.favoriteRepo, req)
}

func (h *contentHandler) GetUserFavorite(ctx context.Context, req *proto.UIDPageReq) (*proto.InfoIDsResp, error) {
	return h._behaviorFindInfoByUser(ctx, h.favoriteRepo, req)
}

func (h *contentHandler) GetInfoFavorite(ctx context.Context, req *proto.InfoIDPageReq) (*proto.UserIDsResp, error) {
	return h._behaviorFindUserByInfo(ctx, h.favoriteRepo, req)
}

func (h *contentHandler) DeleteFavorite(ctx context.Context, req *proto.InfoIDAndUIDReq) (*empty.Empty, error) {
	return h._behaviorDelete(ctx, h.favoriteRepo, req)
}
