package models

import "time"

type Tag struct {
	Tag         string    `bson:"_id"`
	Hot         uint64    `bson:"hot"`
	CreateTime  time.Time `bson:"createTime"`
	LastUseTime time.Time `bson:"lastUseTime"`
}
