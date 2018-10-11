package account

import (
	"context"
	"errors"
	"github.com/golang/protobuf/ptypes/empty"
	log "github.com/sirupsen/logrus"
	"github.com/zhsyourai/teddy-backend/common/models"
	"github.com/zhsyourai/teddy-backend/uaa/components"
	"github.com/zhsyourai/teddy-backend/uaa/converter"
	"github.com/zhsyourai/teddy-backend/uaa/proto"
	"github.com/zhsyourai/teddy-backend/uaa/repositories"
	"golang.org/x/crypto/bcrypt"
	"time"
)

var UserNotFoundErr = errors.New("user not found")
var UserHasBeenRegisteredErr = errors.New("user has been registered")
var OldPasswordNotCorrectErr = errors.New("old password not correct")
var PasswordModifyErr = errors.New("password modify error")

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

func (h *accountHandler) GetAll(context.Context, *empty.Empty) (*proto.GetAllResp, error) {
	accs, err := h.repo.FindAll()
	if err != nil {
		log.Error(err)
		return nil, err
	}

	var resp proto.GetAllResp
	for _, v := range accs {
		var pbacc proto.Account
		converter.CopyFromAccountToPBAccount(&v, &pbacc)
		resp.Accounts = append(resp.Accounts, &pbacc)
	}

	return &resp, nil
}

func (h *accountHandler) GetByUsername(ctx context.Context, req *proto.GetByUsernameReq) (*proto.Account, error) {
	if err := validateGetByUsernameReq(req); err != nil {
		return nil, err
	}

	acc, err := h.repo.FindAccountByUsername(req.GetUsername())
	if err != nil {
		log.Error(err)
		return nil, err
	}
	var resp proto.Account
	converter.CopyFromAccountToPBAccount(&acc, &resp)
	return &resp, nil
}

func (h *accountHandler) DeleteByUsername(ctx context.Context, req *proto.DeleteByUsernameReq) (*empty.Empty, error) {
	if err := validateDeleteByUsernameReq(req); err != nil {
		return nil, err
	}

	err := h.repo.DeleteAccountByUsername(req.GetUsername())
	if err != nil {
		log.Error(err)
		return nil, err
	}
	var resp empty.Empty
	return &resp, nil
}

func (h *accountHandler) Register(ctx context.Context, req *proto.RegisterReq) (*proto.Account, error) {
	if err := validateRegisterReq(req); err != nil {
		return nil, err
	}

	_, err := h.repo.FindAccountByUsername(req.GetUsername())
	if err == nil {
		return nil, err
	}
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.GetPassword()), bcrypt.DefaultCost)
	if err != nil {
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
	// TODO: Check Roles
	account.Roles = req.GetRoles()
	account.CreateDate = time.Now()
	account.Email = req.Email
	account.CredentialsExpired = false
	account.AccountLocked = false
	account.AccountExpired = false

	err = h.repo.InsertAccount(&account)
	if err != nil {
		log.Error(err)
		return nil, err
	}
	var resp proto.Account
	converter.CopyFromAccountToPBAccount(&account, &resp)
	return &resp, nil
}

func (h *accountHandler) VerifyPassword(ctx context.Context, req *proto.VerifyPasswordReq) (*proto.Account, error) {
	if err := validateVerifyPasswordReq(req); err != nil {
		return nil, err
	}

	acc, err := h.repo.FindAccountByUsername(req.GetUsername())
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
	converter.CopyFromAccountToPBAccount(&acc, &resp)
	return &resp, nil
}

func (h *accountHandler) ChangePassword(ctx context.Context, req *proto.ChangePasswordReq) (*empty.Empty, error) {
	if err := validateChangePasswordReq(req); err != nil {
		return nil, err
	}

	acc, err := h.repo.FindAccountByUsername(req.GetUsername())
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
	err = h.repo.UpdateAccountByUsername(req.GetUsername(), map[string]interface{}{
		"password": hashedPassword,
	})
	if err != nil {
		log.Error(err)
		return nil, PasswordModifyErr
	}
	var resp empty.Empty
	return &resp, nil
}
