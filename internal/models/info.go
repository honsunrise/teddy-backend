package models

import (
	"github.com/mongodb/mongo-go-driver/bson/primitive"
	"time"
)

type TypeAndTag struct {
	Type string `bson:"type"`
	Tag  string `bson:"tag"`
}

type Info struct {
	ID               primitive.ObjectID `bson:"_id"`
	UID              string             `bson:"uid"`
	Author           string             `bson:"author"`
	Title            string             `bson:"title"`
	Summary          string             `bson:"summary"`
	Country          string             `bson:"country"`
	ContentTime      time.Time          `bson:"contentTime"`
	CoverResources   map[string]string  `bson:"coverResources"`
	PublishTime      time.Time          `bson:"publishTime"`
	LastReviewTime   time.Time          `bson:"lastReviewTime"`
	Valid            bool               `bson:"valid"`
	WatchCount       uint64             `bson:"watchCount"`
	Tags             []*TypeAndTag      `bson:"tags"`
	LatestModifyTime time.Time          `bson:"lastModifyTime"`
	CanReview        bool               `bson:"canReview"`
	Archived         bool               `bson:"archived"`
	LatestSegmentID  primitive.ObjectID `bson:"latestSegmentID"`
	SegmentCount     uint64             `bson:"segmentCount"`
}
