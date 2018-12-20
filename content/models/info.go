package models

import (
	"github.com/mongodb/mongo-go-driver/bson/objectid"
	"time"
)

type InfoType uint32

type Info struct {
	Id             objectid.ObjectID `bson:"_id"`
	UID            string            `bson:"uid"`
	Title          string            `bson:"title"`
	Content        string            `bson:"content"`
	ContentTime    time.Time         `bson:"contentTime"`
	CoverResources map[string]string `bson:"coverResources"`
	PublishTime    time.Time         `bson:"publishTime"`
	LastReviewTime time.Time         `bson:"lastReviewTime"`
	Valid          bool              `bson:"valid"`
	WatchCount     int64             `bson:"watchCount"`
	Tags           []string          `bson:"tags"`
	LastModifyTime time.Time         `bson:"lastModifyTime"`
	CanReview      bool              `bson:"canReview"`
	ThumbUp        int64             `bson:"thumbUp"`
	ThumbDown      int64             `bson:"thumbDown"`
	Favorites      int64             `bson:"favorites"`
}
