package server

import (
	"github.com/golang/protobuf/ptypes/empty"
	"github.com/zhsyourai/teddy-backend/content/proto"
	"github.com/zhsyourai/teddy-backend/content/repositories"
	context2 "golang.org/x/net/context"
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

func (h *contentHandler) GetTags(context2.Context, *proto.GetTagReq) (*proto.GetTagsResp, error) {
	panic("implement me")
}

func (h *contentHandler) PublishInfo(context2.Context, *proto.PublishInfoReq) (*empty.Empty, error) {
	panic("implement me")
}

func (h *contentHandler) EditInfo(context2.Context, *proto.EditInfoReq) (*empty.Empty, error) {
	panic("implement me")
}

func (h *contentHandler) GetInfos(context2.Context, *proto.GetInfosReq) (*proto.GetInfosResp, error) {
	panic("implement me")
}

func (h *contentHandler) DeleteInfo(context2.Context, *proto.InfoIdReq) (*empty.Empty, error) {
	panic("implement me")
}

func (h *contentHandler) WatchInfo(context2.Context, *proto.InfoIdReq) (*empty.Empty, error) {
	panic("implement me")
}

func (h *contentHandler) LikeInfo(context2.Context, *proto.InfoIdReq) (*empty.Empty, error) {
	panic("implement me")
}

func (h *contentHandler) UnLikeInfo(context2.Context, *proto.InfoIdReq) (*empty.Empty, error) {
	panic("implement me")
}

func (h *contentHandler) GetUserLikes(context2.Context, *proto.UidPageReq) (*proto.InfoIdsResp, error) {
	panic("implement me")
}

func (h *contentHandler) GetUserUnlikes(context2.Context, *proto.UidPageReq) (*proto.InfoIdsResp, error) {
	panic("implement me")
}

func (h *contentHandler) GetInfoLiked(context2.Context, *proto.InfoIdPageReq) (*proto.UserIdsResp, error) {
	panic("implement me")
}

func (h *contentHandler) GetInfoUnliked(context2.Context, *proto.InfoIdPageReq) (*proto.UserIdsResp, error) {
	panic("implement me")
}

func (h *contentHandler) FavoriteInfo(context2.Context, *proto.InfoIdReq) (*empty.Empty, error) {
	panic("implement me")
}

func (h *contentHandler) GetUserFavorite(context2.Context, *proto.UidPageReq) (*proto.InfoIdsResp, error) {
	panic("implement me")
}

func (h *contentHandler) GetInfoFavorited(context2.Context, *proto.InfoIdPageReq) (*proto.UserIdsResp, error) {
	panic("implement me")
}
