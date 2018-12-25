package models

import "github.com/mongodb/mongo-go-driver/bson/objectid"

type Segment struct {
	ID         objectid.ObjectID `bson:"_id"`
	InfoID     objectid.ObjectID `bson:"infoID"`
	Title      string            `bson:"title"`
	No         int64             `bson:"no"`
	Labels     []string          `bson:"labels"`
	Content    map[string]string `bson:"content"`
	WatchCount int64             `bson:"watchCount"`
}
