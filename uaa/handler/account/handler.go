package account

import (
	"context"
	"errors"
	"github.com/casbin/casbin"
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

func NewAccountHandler(repo repositories.AccountRepository,
	uidGen components.UidGenerator, enforcer *casbin.Enforcer) (proto.UAAHandler, error) {
	return &accountHandler{
		repo:     repo,
		uidGen:   uidGen,
		enforcer: enforcer,
	}, nil
}

type accountHandler struct {
	repo     repositories.AccountRepository
	uidGen   components.UidGenerator
	enforcer *casbin.Enforcer
}

func (h *accountHandler) GetAll(ctx context.Context, req *empty.Empty, resp *proto.GetAllResp) error {
	accs, err := h.repo.FindAll()
	if err != nil {
		log.Error(err)
		return err
	}

	for _, v := range accs {
		var pbacc proto.Account
		converter.CopyFromAccountToPBAccount(&v, &pbacc)
		resp.Accounts = append(resp.Accounts, &pbacc)
	}

	return nil
}

func (h *accountHandler) GetByUsername(ctx context.Context, req *proto.GetByUsernameReq, resp *proto.Account) error {
	if err := validateGetByUsernameReq(req); err != nil {
		return err
	}

	acc, err := h.repo.FindAccountByUsername(req.GetUsername())
	if err != nil {
		log.Error(err)
		return UserNotFoundErr
	}
	converter.CopyFromAccountToPBAccount(&acc, resp)
	return nil
}

func (h *accountHandler) DeleteByUsername(ctx context.Context, req *proto.DeleteByUsernameReq, resp *empty.Empty) error {
	if err := validateDeleteByUsernameReq(req); err != nil {
		return err
	}

	err := h.repo.DeleteAccountByUsername(req.GetUsername())
	if err != nil {
		log.Error(err)
		return UserNotFoundErr
	}
	return nil
}

func (h *accountHandler) Register(ctx context.Context, req *proto.RegisterReq, resp *proto.Account) error {
	if err := validateRegisterReq(req); err != nil {
		return err
	}

	_, err := h.repo.FindAccountByUsername(req.GetUsername())
	if err == nil {
		return UserHasBeenRegisteredErr
	}
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.GetPassword()), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	uid, err := h.uidGen.NexID()
	if err != nil {
		log.Error(err)
		return err
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
		return err
	}
	converter.CopyFromAccountToPBAccount(&account, resp)
	return nil
}

func (h *accountHandler) VerifyPassword(ctx context.Context, req *proto.VerifyPasswordReq, resp *proto.Account) error {
	if err := validateVerifyPasswordReq(req); err != nil {
		return err
	}

	acc, err := h.repo.FindAccountByUsername(req.GetUsername())
	if err != nil {
		log.Error(err)
		return UserNotFoundErr
	}
	err = bcrypt.CompareHashAndPassword(acc.Password, []byte(req.GetPassword()))
	if err != nil {
		log.Error(err)
		return UserNotFoundErr
	}
	converter.CopyFromAccountToPBAccount(&acc, resp)
	return nil
}

func (h *accountHandler) ChangePassword(ctx context.Context, req *proto.ChangePasswordReq, resp *empty.Empty) error {
	if err := validateChangePasswordReq(req); err != nil {
		return err
	}

	acc, err := h.repo.FindAccountByUsername(req.GetUsername())
	if err != nil {
		log.Error(err)
		return UserNotFoundErr
	}
	err = bcrypt.CompareHashAndPassword(acc.Password, []byte(req.GetOldPassword()))
	if err != nil {
		log.Error(err)
		return OldPasswordNotCorrectErr
	}
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.GetNewPassword()), bcrypt.DefaultCost)
	if err != nil {
		log.Error(err)
		return PasswordModifyErr
	}
	err = h.repo.UpdateAccountByUsername(req.GetUsername(), map[string]interface{}{
		"password": hashedPassword,
	})
	if err != nil {
		log.Error(err)
		return PasswordModifyErr
	}
	return nil
}
