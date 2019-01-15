package uaa

import (
	"context"
	"errors"
	"github.com/golang/protobuf/ptypes"
	"github.com/golang/protobuf/ptypes/empty"
	"github.com/mongodb/mongo-go-driver/mongo"
	log "github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
	"teddy-backend/internal/components"
	"teddy-backend/internal/models"
	"teddy-backend/internal/proto/uaa"
	"teddy-backend/internal/repositories"
	"time"
)

func NewAccountServer(repo repositories.AccountRepository, uidGen components.UidGenerator) (uaa.UAAServer, error) {
	return &accountHandler{
		repo:   repo,
		uidGen: uidGen,
	}, nil
}

type accountHandler struct {
	repo   repositories.AccountRepository
	uidGen components.UidGenerator
}

func (h *accountHandler) GetAll(ctx context.Context, req *uaa.GetAllReq) (*uaa.GetAllResp, error) {
	accounts, err := h.repo.FindAll(req.Page, req.Size, req.Sorts)
	if err != nil {
		log.Error(err)
		return nil, err
	}

	var resp uaa.GetAllResp
	for _, v := range accounts {
		var pbAcc uaa.Account
		copyFromAccountToPBAccount(v, &pbAcc)
		resp.Accounts = append(resp.Accounts, &pbAcc)
	}

	return &resp, nil
}

func (h *accountHandler) GetOne(ctx context.Context, req *uaa.GetOneReq) (*uaa.Account, error) {
	if err := validateGetOneReq(req); err != nil {
		return nil, err
	}

	acc, err := h.repo.FindOne(req.GetPrincipal())
	if err != nil {
		log.Error(err)
		return nil, err
	}
	var resp uaa.Account
	copyFromAccountToPBAccount(acc, &resp)
	return &resp, nil
}

func (h *accountHandler) RegisterByNormal(ctx context.Context, req *uaa.RegisterNormalReq) (*uaa.Account, error) {
	if err := validateRegisterNormalReq(req); err != nil {
		return nil, err
	}

	_, err := h.repo.FindOne(req.GetUsername())
	if err != mongo.ErrNoDocuments {
		return nil, ErrAccountExist
	}

	var account models.Account
	account.Username = req.GetUsername()
	account.Roles = req.GetRoles()
	account.CreateDate = time.Now()
	account.OAuthUIds = make(map[string]string)
	account.CredentialsExpired = false
	account.Locked = false

	if x, ok := req.GetContact().(*uaa.RegisterNormalReq_Email); ok {
		_, err = h.repo.FindOne(req.GetEmail())
		if err != mongo.ErrNoDocuments {
			return nil, ErrAccountExist
		}
		account.Email = x.Email
	} else if x, ok := req.GetContact().(*uaa.RegisterNormalReq_Phone); ok {
		_, err = h.repo.FindOne(req.GetPhone())
		if err != mongo.ErrNoDocuments {
			return nil, ErrAccountExist
		}
		account.Phone = x.Phone
	} else {
		panic(errors.New("never happen"))
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.GetPassword()), bcrypt.DefaultCost)
	if err != nil {
		log.Error(err)
		return nil, err
	}
	account.Password = hashedPassword

	uid, err := h.uidGen.NexID()
	if err != nil {
		log.Error(err)
		return nil, err
	}
	account.UID = uid

	err = h.repo.InsertAccount(&account)
	if err != nil {
		log.Error(err)
		return nil, err
	}
	var resp uaa.Account
	copyFromAccountToPBAccount(&account, &resp)
	return &resp, nil
}

func (h *accountHandler) RegisterByOAuth(ctx context.Context, req *uaa.RegisterOAuthReq) (*uaa.Account, error) {
	panic("implement me")
}

func (h *accountHandler) VerifyPassword(ctx context.Context, req *uaa.VerifyAccountReq) (*uaa.Account, error) {
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
	var resp uaa.Account
	copyFromAccountToPBAccount(acc, &resp)
	return &resp, nil
}

func (h *accountHandler) DeleteOne(ctx context.Context, req *uaa.UIDReq) (*empty.Empty, error) {
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

func (h *accountHandler) DoLockAccount(ctx context.Context, req *uaa.UIDReq) (*empty.Empty, error) {
	panic("implement me")
}

func (h *accountHandler) DoCredentialsExpired(ctx context.Context, req *uaa.UIDReq) (*empty.Empty, error) {
	panic("implement me")
}

func (h *accountHandler) UpdateSignIn(ctx context.Context, req *uaa.UpdateSignInReq) (*empty.Empty, error) {
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

func (h *accountHandler) ChangePassword(ctx context.Context, req *uaa.ChangePasswordReq) (*empty.Empty, error) {
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
