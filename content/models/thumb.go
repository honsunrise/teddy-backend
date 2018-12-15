package models

import (
	"github.com/mongodb/mongo-go-driver/bson/objectid"
	"time"
)

type ThumbItem struct {
	InfoId       objectid.ObjectID `bson:"infoId"`
	FavoriteTime time.Time         `bson:"time"`
}

type Thumb struct {
	Id        objectid.ObjectID `bson:"_id"`
	UserId    objectid.ObjectID `bson:"userId"`
	FirstTime time.Time         `bson:"firstTime"`
	LastTime  time.Time         `bson:"lastTime"`
	Count     uint64            `bson:"count"`
	Items     []ThumbItem       `bson:"items"`
}
