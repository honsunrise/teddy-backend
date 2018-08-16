package models

import "time"

type OAuthType int

const (
	GOOGLE OAuthType = iota
	WECHAT
	QQ
	SINA
	FACEBOOK
	TWITTER
)

type Account struct {
	UID                string               `bson:"_id" json:"uid"`
	Username           string               `bson:"username" json:"username"`
	Email              string               `bson:"email" json:"email"`
	CreateDate         time.Time            `bson:"create_date" json:"create_date"`
	OAuthUserIds       map[OAuthType]string `bson:"oauth_user_ids" json:"oauth_user_ids"`
	Password           []byte               `bson:"password" json:"password"`
	AccountExpired     bool                 `bson:"account_expired" json:"account_expired"`
	AccountLocked      bool                 `bson:"account_locked" json:"account_locked"`
	CredentialsExpired bool                 `bson:"credentials_expired" json:"credentials_expired"`
	Roles              []string             `bson:"roles" json:"roles"`
	UpdateDate         time.Time            `bson:"update_date" json:"update_date"`
}
