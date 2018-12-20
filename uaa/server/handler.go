package server

import (
	"context"
	"errors"
	"github.com/golang/protobuf/ptypes"
	"github.com/golang/protobuf/ptypes/empty"
	"github.com/mongodb/mongo-go-driver/mongo"
	log "github.com/sirupsen/logrus"
	"github.com/zhsyourai/teddy-backend/common/proto"
	"github.com/zhsyourai/teddy-backend/uaa/components"
	"github.com/zhsyourai/teddy-backend/uaa/converter"
	"github.com/zhsyourai/teddy-backend/uaa/models"
	"github.com/zhsyourai/teddy-backend/uaa/repositories"
	"golang.org/x/crypto/bcrypt"
	"time"
)

func NewAccountServer(repo repositories.AccountRepository, uidGen components.UidGenerator) (proto.UAAServer, error) {
	return &accountHandler{
		repo:   repo,
		uidGen: uidGen,
	}, nil
}

type accountHandler struct {
	repo   repositories.AccountRepository
	uidGen components.UidGenerator
}

func (h *accountHandler) GetAll(ctx context.Context, req *proto.GetAllReq) (*proto.GetAllResp, error) {
	accounts, err := h.repo.FindAll(req.Page, req.Size, req.Sorts)
	if err != nil {
		log.Error(err)
		return nil, err
	}

	var resp proto.GetAllResp
	for _, v := range accounts {
		var pbAcc proto.Account
		converter.CopyFromAccountToPBAccount(v, &pbAcc)
		resp.Accounts = append(resp.Accounts, &pbAcc)
	}

	return &resp, nil
}

func (h *accountHandler) GetOne(ctx context.Context, req *proto.GetOneReq) (*proto.Account, error) {
	if err := validateGetOneReq(req); err != nil {
		return nil, err
	}

	acc, err := h.repo.FindOne(req.GetPrincipal())
	if err != nil {
		log.Error(err)
		return nil, err
	}
	var resp proto.Account
	converter.CopyFromAccountToPBAccount(acc, &resp)
	return &resp, nil
}

func (h *accountHandler) RegisterByNormal(ctx context.Context, req *proto.RegisterNormalReq) (*proto.Account, error) {
	if err := validateRegisterNormalReq(req); err != nil {
		return nil, err
	}

	_, err := h.repo.FindOne(req.GetUsername())
	if err != mongo.ErrNoDocuments {
		return nil, ErrAccountExist
	}

	_, err = h.repo.FindOne(req.GetEmail())
	if err != mongo.ErrNoDocuments {
		return nil, ErrAccountExist
	}

	_, err = h.repo.FindOne(req.GetPhone())
	if err != mongo.ErrNoDocuments {
		return nil, ErrAccountExist
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.GetPassword()), bcrypt.DefaultCost)
	if err != nil {
		log.Error(err)
		return nil, err
	}
	uid, err := h.uidGen.NexID()
	if err != nil {
		log.Error(err)
		return nil, err
	}

	var account models.Account
	account.UID = uid
	account.Username = req.GetUsername()
	account.Password = hashedPassword
	account.Roles = req.GetRoles()
	account.CreateDate = time.Now()
	account.OAuthUIds = make(map[string]string)
	account.CredentialsExpired = false
	account.Locked = false

	if x, ok := req.GetContact().(*proto.RegisterNormalReq_Email); ok {
		account.Email = x.Email
	} else if x, ok := req.GetContact().(*proto.RegisterNormalReq_Phone); ok {
		account.Phone = x.Phone
	} else {
		panic(errors.New("never happen"))
	}

	err = h.repo.InsertAccount(&account)
	if err != nil {
		log.Error(err)
		return nil, err
	}
	var resp proto.Account
	converter.CopyFromAccountToPBAccount(&account, &resp)
	return &resp, nil
}

func (h *accountHandler) RegisterByOAuth(ctx context.Context, req *proto.RegisterOAuthReq) (*proto.Account, error) {
	panic("implement me")
}

func (h *accountHandler) VerifyPassword(ctx context.Context, req *proto.VerifyAccountReq) (*proto.Account, error) {
	if err := validateVerifyPasswordReq(req); err != nil {
		return nil, err
	}

	acc, err := h.repo.FindOne(req.GetPrincipal())
	if err != nil {
		log.Error(err)
		return nil, UserNotFoundErr
	}
	err = bcrypt.CompareHashAndPassword(acc.Password, []byte(req.GetPassword()))
	if err != nil {
		log.Error(err)
		return nil, UserNotFoundErr
	}
	var resp proto.Account
	converter.CopyFromAccountToPBAccount(acc, &resp)
	return &resp, nil
}

func (h *accountHandler) DeleteOne(ctx context.Context, req *proto.UIDReq) (*empty.Empty, error) {
	if err := validateUIDReq(req); err != nil {
		return nil, err
	}

	err := h.repo.DeleteOne(req.GetUid())
	if err != nil {
		log.Error(err)
		return nil, err
	}
	var resp empty.Empty
	return &resp, nil
}

func (h *accountHandler) DoLockAccount(ctx context.Context, req *proto.UIDReq) (*empty.Empty, error) {
	panic("implement me")
}

func (h *accountHandler) DoCredentialsExpired(ctx context.Context, req *proto.UIDReq) (*empty.Empty, error) {
	panic("implement me")
}

func (h *accountHandler) UpdateSignIn(ctx context.Context, req *proto.UpdateSignInReq) (*empty.Empty, error) {
	if err := validateUpdateSignInReq(req); err != nil {
		return nil, err
	}

	acc, err := h.repo.FindOne(req.GetPrincipal())
	if err != nil {
		log.Error(err)
		return nil, UserNotFoundErr
	}

	lastSignInTime, err := ptypes.Timestamp(req.Time)
	if err != nil {
		log.Error(err)
		return nil, err
	}

	err = h.repo.UpdateOne(acc.UID, map[string]interface{}{
		"last_sign_in_ip":   req.Ip,
		"last_sign_in_time": lastSignInTime,
	})
	if err != nil {
		log.Error(err)
		return nil, PasswordModifyErr
	}
	var resp empty.Empty
	return &resp, nil
}

func (h *accountHandler) ChangePassword(ctx context.Context, req *proto.ChangePasswordReq) (*empty.Empty, error) {
	if err := validateChangePasswordReq(req); err != nil {
		return nil, err
	}

	acc, err := h.repo.FindOne(req.GetPrincipal())
	if err != nil {
		log.Error(err)
		return nil, UserNotFoundErr
	}
	err = bcrypt.CompareHashAndPassword(acc.Password, []byte(req.GetOldPassword()))
	if err != nil {
		log.Error(err)
		return nil, OldPasswordNotCorrectErr
	}
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.GetNewPassword()), bcrypt.DefaultCost)
	if err != nil {
		log.Error(err)
		return nil, PasswordModifyErr
	}
	err = h.repo.UpdateOne(acc.UID, map[string]interface{}{
		"password": hashedPassword,
	})
	if err != nil {
		log.Error(err)
		return nil, PasswordModifyErr
	}
	var resp empty.Empty
	return &resp, nil
}
