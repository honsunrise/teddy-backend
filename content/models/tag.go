package models

import "time"

type Tag struct {
	Tag         string    `bson:"_id"`
	Usage       uint64    `bson:"usage"`
	CreateTime  time.Time `bson:"createTime"`
	LastUseTime time.Time `bson:"lastUseTime"`
}
