package models

import (
	"github.com/mongodb/mongo-go-driver/bson/primitive"
	"time"
)

type Tag struct {
	ID          primitive.ObjectID `bson:"_id"`
	Tag         string             `bson:"tag"`
	Type        string             `bson:"type"`
	Usage       uint64             `bson:"usage"`
	CreateTime  time.Time          `bson:"createTime"`
	LastUseTime time.Time          `bson:"lastUseTime"`
}
