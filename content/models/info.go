package models

import (
	"time"
)

type Info struct {
	Id             string    `bson:"_id"`
	Uid            string    `bson:"uid"`
	Title          string    `bson:"title"`
	Type           uint32    `bson:"type"`
	Content        string    `bson:"content"`
	CoverList      []string  `bson:"coverList"`
	CoverVideo     string    `bson:"coverVideo"`
	PublishTime    time.Time `bson:"publishTime"`
	LastReviewTime time.Time `bson:"lastReviewTime"`
	Valid          bool      `bson:"valid"`
	WatchCount     int64     `bson:"watchCount"`
	Tags           []string  `bson:"tags"`
	Likes          int64     `bson:"likes"`
	IsLike         bool      `bson:"isLike"`
	LikeList       []string  `bson:"likeList"`
	Unlike         int64     `bson:"unlike"`
	IsUnlike       bool      `bson:"isUnlike"`
	UnlikeList     []string  `bson:"unlikeList"`
	Favorites      int64     `bson:"favorites"`
	IsFavorite     bool      `bson:"isFavorite"`
	FavoriteList   []string  `bson:"favoriteList"`
	LastModifyTime time.Time `bson:"lastModifyTime"`
	CanReview      bool      `bson:"canReview"`
}
