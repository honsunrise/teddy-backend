package converter

import (
	"github.com/golang/protobuf/ptypes/timestamp"
	"github.com/zhsyourai/teddy-backend/common/proto"
	"github.com/zhsyourai/teddy-backend/uaa/models"
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
	pbacc.OauthUserIds = acc.OAuthUserIds
	pbacc.CreateDate = &timestamp.Timestamp{
		Seconds: acc.CreateDate.Unix(),
		Nanos:   int32(acc.CreateDate.Nanosecond()),
	}
	pbacc.UpdateDate = &timestamp.Timestamp{
		Seconds: acc.UpdateDate.Unix(),
		Nanos:   int32(acc.UpdateDate.Nanosecond()),
	}
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
	acc.OAuthUserIds = pbacc.OauthUserIds
	acc.CreateDate = time.Unix(pbacc.CreateDate.Seconds, int64(pbacc.CreateDate.Nanos))
	acc.UpdateDate = time.Unix(pbacc.UpdateDate.Seconds, int64(pbacc.UpdateDate.Nanos))
}
