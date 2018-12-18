package models

import (
	"time"
)

const (
	GOOGLE   string = "GOOGLE"
	WECHAT   string = "WECHAT"
	QQ       string = "QQ"
	SINA     string = "SINA"
	FACEBOOK string = "FACEBOOK"
	TWITTER  string = "TWITTER"
)

type Account struct {
	UID                string            `bson:"_id" json:"uid"`
	Username           string            `bson:"username" json:"username"`
	Email              string            `bson:"email" json:"email"`
	Phone              string            `bson:"phone" json:"phone"`
	CreateDate         time.Time         `bson:"create_date" json:"create_date"`
	OAuthUserIds       map[string]string `bson:"oauth_user_ids" json:"oauth_user_ids"`
	Password           []byte            `bson:"password" json:"password"`
	AccountExpired     bool              `bson:"account_expired" json:"account_expired"`
	AccountLocked      bool              `bson:"account_locked" json:"account_locked"`
	CredentialsExpired bool              `bson:"credentials_expired" json:"credentials_expired"`
	Roles              []string          `bson:"roles" json:"roles"`
	UpdateDate         time.Time         `bson:"update_date" json:"update_date"`
}
