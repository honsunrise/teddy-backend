package server

import (
	"github.com/golang/protobuf/ptypes/empty"
	"github.com/zhsyourai/teddy-backend/common/proto"
	"github.com/zhsyourai/teddy-backend/content/repositories"
	"golang.org/x/net/context"
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

func (h *contentHandler) GetTags(ctx context.Context, req *proto.GetTagReq) (*proto.GetTagsResp, error) {
	panic("implement me")
}

func (h *contentHandler) PublishInfo(ctx context.Context, req *proto.PublishInfoReq) (*empty.Empty, error) {
	panic("implement me")
}

func (h *contentHandler) EditInfo(ctx context.Context, req *proto.EditInfoReq) (*empty.Empty, error) {
	panic("implement me")
}

func (h *contentHandler) GetInfos(ctx context.Context, req *proto.GetInfosReq) (*proto.GetInfosResp, error) {
	panic("implement me")
}

func (h *contentHandler) DeleteInfo(ctx context.Context, req *proto.InfoIdReq) (*empty.Empty, error) {
	panic("implement me")
}

func (h *contentHandler) WatchInfo(ctx context.Context, req *proto.InfoIdReq) (*empty.Empty, error) {
	panic("implement me")
}

func (h *contentHandler) LikeInfo(ctx context.Context, req *proto.InfoIdReq) (*empty.Empty, error) {
	panic("implement me")
}

func (h *contentHandler) UnLikeInfo(ctx context.Context, req *proto.InfoIdReq) (*empty.Empty, error) {
	panic("implement me")
}

func (h *contentHandler) GetUserLikes(ctx context.Context, req *proto.UidPageReq) (*proto.InfoIdsResp, error) {
	panic("implement me")
}

func (h *contentHandler) GetUserUnlikes(ctx context.Context, req *proto.UidPageReq) (*proto.InfoIdsResp, error) {
	panic("implement me")
}

func (h *contentHandler) GetInfoLiked(ctx context.Context, req *proto.InfoIdPageReq) (*proto.UserIdsResp, error) {
	panic("implement me")
}

func (h *contentHandler) GetInfoUnliked(ctx context.Context, req *proto.InfoIdPageReq) (*proto.UserIdsResp, error) {
	panic("implement me")
}

func (h *contentHandler) FavoriteInfo(ctx context.Context, req *proto.InfoIdReq) (*empty.Empty, error) {
	panic("implement me")
}

func (h *contentHandler) GetUserFavorite(ctx context.Context, req *proto.UidPageReq) (*proto.InfoIdsResp, error) {
	panic("implement me")
}

func (h *contentHandler) GetInfoFavorited(ctx context.Context, req *proto.InfoIdPageReq) (*proto.UserIdsResp, error) {
	panic("implement me")
}
