package converter

import (
	"github.com/zhsyourai/teddy-backend/common/models"
	"github.com/zhsyourai/teddy-backend/uaa/proto"
	"time"
)

func CopyFromAccountToPBAccount(acc *models.Account, pbacc *proto.Account) {
	if acc == nil || pbacc == nil {
		return
	}
	pbacc.Uid = acc.UID
	pbacc.Username = acc.Username
	pbacc.Email = acc.Email
	pbacc.Password = acc.Password
	pbacc.AccountExpired = acc.AccountExpired
	pbacc.AccountLocked = acc.AccountLocked
	pbacc.CredentialsExpired = acc.CredentialsExpired
	pbacc.Roles = acc.Roles
	for k, v := range acc.OAuthUserIds {
		pbacc.OauthUserIds[uint32(k)] = v
	}
	pbacc.CreateDate.Seconds = acc.CreateDate.Unix()
	pbacc.CreateDate.Nanos = int32(acc.CreateDate.Nanosecond())

	pbacc.UpdateDate.Seconds = acc.UpdateDate.Unix()
	pbacc.UpdateDate.Nanos = int32(acc.UpdateDate.Nanosecond())
}

func CopyFromPBAccountToAccount(pbacc *proto.Account, acc *models.Account) {
	if acc == nil || pbacc == nil {
		return
	}
	acc.UID = pbacc.Uid
	acc.Username = pbacc.Username
	acc.Email = pbacc.Email
	acc.Password = pbacc.Password
	acc.AccountExpired = pbacc.AccountExpired
	acc.AccountLocked = pbacc.AccountLocked
	acc.CredentialsExpired = pbacc.CredentialsExpired
	acc.Roles = pbacc.Roles
	for k, v := range pbacc.OauthUserIds {
		acc.OAuthUserIds[models.OAuthType(k)] = v
	}
	acc.CreateDate = time.Unix(pbacc.CreateDate.Seconds, int64(pbacc.CreateDate.Nanos))
	acc.UpdateDate = time.Unix(pbacc.UpdateDate.Seconds, int64(pbacc.UpdateDate.Nanos))
}
