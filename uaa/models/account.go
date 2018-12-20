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
	UID                string            `bson:"_id"`
	Username           string            `bson:"username"`
	Email              string            `bson:"email"`
	Phone              string            `bson:"phone"`
	OAuthUIds          map[string]string `bson:"oauth_uids"`
	Roles              []string          `bson:"roles"`
	Password           []byte            `bson:"password"`
	Locked             bool              `bson:"locked"`
	CredentialsExpired bool              `bson:"credentials_expired"`
	CreateDate         time.Time         `bson:"create_date"`
	UpdateDate         time.Time         `bson:"update_date"`
	LastSignInIP       string            `bson:"last_sign_in_ip"`
	LastSignInTime     time.Time         `bson:"last_sign_in_time"`
}
