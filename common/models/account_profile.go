package models

import "time"

type GenderType int

const (
	MAN GenderType = iota
	WOMAN
	UNKNOWN
)

type AccountProfile struct {
	UID        string     `bson:"_id" json:"uid"`
	Firstname  string     `bson:"firstname" json:"firstname"`
	Lastname   string     `bson:"lastname" json:"lastname"`
	AvatarUrl  string     `bson:"avatar_url" json:"avatar_url"`
	Bio        string     `bson:"bio" json:"bio"`
	Birthday   time.Time  `bson:"birthday" json:"birthday"`
	Gender     GenderType `bson:"gender" json:"gender"`
	UpdateDate time.Time  `bson:"update_date" json:"update_date"`
	Locale     string     `bson:"locale" json:"locale"`
}
