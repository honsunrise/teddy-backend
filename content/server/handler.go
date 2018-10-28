package server

import (
	"context"
	"github.com/golang/protobuf/ptypes/empty"
	"github.com/zhsyourai/teddy-backend/content/proto"
	"github.com/zhsyourai/teddy-backend/content/repositories"
	"gopkg.in/gomail.v2"
)

func NewContentServer(repo repositories.InfoRepository) (proto.ContentServer, error) {
	instance := &contentHandler{
		repo:    repo,
		mailCh:  make(chan *gomail.Message),
		mailErr: make(chan error),
	}
	return instance, nil
}

type contentHandler struct {
	repo    repositories.InfoRepository
	mailCh  chan *gomail.Message
	mailErr chan error
}

func (h *contentHandler) GetTags(context.Context, *proto.GetTagReq, *proto.GetTagsResp) error {
	panic("implement me")
}

func (h *contentHandler) PublishInfo(context.Context, *proto.PublishInfoReq, *empty.Empty) error {
	panic("implement me")
}

func (h *contentHandler) EditInfo(context.Context, *proto.EditInfoReq, *empty.Empty) error {
	panic("implement me")
}

func (h *contentHandler) GetInfos(context.Context, *proto.GetInfosReq, *proto.GetInfosResp) error {
	panic("implement me")
}

func (h *contentHandler) DeleteInfo(context.Context, *proto.InfoIdReq, *empty.Empty) error {
	panic("implement me")
}

func (h *contentHandler) WatchInfo(context.Context, *proto.InfoIdReq, *empty.Empty) error {
	panic("implement me")
}

func (h *contentHandler) LikeInfo(context.Context, *proto.InfoIdReq, *empty.Empty) error {
	panic("implement me")
}

func (h *contentHandler) UnLikeInfo(context.Context, *proto.InfoIdReq, *empty.Empty) error {
	panic("implement me")
}

func (h *contentHandler) GetUserLikes(context.Context, *proto.UidPageReq, *proto.InfoIdsResp) error {
	panic("implement me")
}

func (h *contentHandler) GetUserUnlikes(context.Context, *proto.UidPageReq, *proto.InfoIdsResp) error {
	panic("implement me")
}

func (h *contentHandler) GetInfoLiked(context.Context, *proto.InfoIdPageReq, *proto.UserIdsResp) error {
	panic("implement me")
}

func (h *contentHandler) GetInfoUnliked(context.Context, *proto.InfoIdPageReq, *proto.UserIdsResp) error {
	panic("implement me")
}

func (h *contentHandler) FavoriteInfo(context.Context, *proto.InfoIdReq, *empty.Empty) error {
	panic("implement me")
}

func (h *contentHandler) GetUserFavorite(context.Context, *proto.UidPageReq, *proto.InfoIdsResp) error {
	panic("implement me")
}

func (h *contentHandler) GetInfoFavorited(context.Context, *proto.InfoIdPageReq, *proto.UserIdsResp) error {
	panic("implement me")
}
