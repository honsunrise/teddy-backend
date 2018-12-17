package server

import (
	"github.com/golang/protobuf/ptypes/empty"
	"github.com/golang/protobuf/ptypes/timestamp"
	"github.com/mongodb/mongo-go-driver/bson/objectid"
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

	for _, tag := range tags {
		pbTag := &proto.Tag{}
		converter.CopyFromTagToPBTag(&tag, pbTag)
		resp.Tags = append(resp.Tags, pbTag)
	}

	return &resp, nil
}

func (h *contentHandler) PublishInfo(ctx context.Context, req *proto.PublishInfoReq) (*empty.Empty, error) {
	var resp empty.Empty
	if err := validatePublishInfoReq(req); err != nil {
		return nil, err
	}

	now := time.Now()
	info := models.Info{
		Id:             objectid.New(),
		UID:            req.Uid,
		Title:          req.Title,
		Content:        req.Content,
		CoverResources: req.CoverResources,
		PublishTime:    now,
		LastReviewTime: nil,
		Valid:          true,
		WatchCount:     0,
		Tags:           req.Tags,
		ThumbUp:        0,
		ThumbUpList:    []string{},
		ThumbDown:      0,
		ThumbDownList:  []string{},
		Favorites:      0,
		FavoriteList:   []string{},
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

	return &resp, nil
}

func (h *contentHandler) GetInfos(ctx context.Context, req *proto.GetInfosReq) (*proto.GetInfosResp, error) {
	var resp proto.GetInfosResp
	if err := validateGetInfosReq(req); err != nil {
		return nil, err
	}

	var infos []models.Info
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

	results := make([]*proto.Info, 0, len(infos))
	for _, info := range infos {
		results = append(results, &proto.Info{
			Id: info.Id.Hex(),
			//TODO: fill other field
		})
	}
	resp.Infos = results
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

	uid, err := objectid.FromHex(req.Uid)
	if err != nil {
		return nil, err
	}

	infoID, err := objectid.FromHex(req.InfoID)
	if err != nil {
		return nil, err
	}

	err = repo.Insert(uid, &models.BehaviorInfoItem{
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

	uid, err := objectid.FromHex(req.Uid)
	if err != nil {
		return nil, err
	}

	items, err := repo.FindInfoByUser(uid, req.Page, req.Size, req.Sorts)

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
			Uid: item.UID.Hex(),
			Time: &timestamp.Timestamp{
				Seconds: item.Time.Unix(),
				Nanos:   int32(item.Time.Nanosecond()),
			},
		})
	}
	resp.Items = results
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

func (h *contentHandler) FavoriteInfo(ctx context.Context, req *proto.InfoIDAndUIDReq) (*empty.Empty, error) {
	return h._behaviorInsert(ctx, h.favoriteRepo, req)
}

func (h *contentHandler) GetUserFavorite(ctx context.Context, req *proto.UIDPageReq) (*proto.InfoIDsResp, error) {
	return h._behaviorFindInfoByUser(ctx, h.favoriteRepo, req)
}

func (h *contentHandler) GetInfoFavorite(ctx context.Context, req *proto.InfoIDPageReq) (*proto.UserIDsResp, error) {
	return h._behaviorFindUserByInfo(ctx, h.favoriteRepo, req)
}
