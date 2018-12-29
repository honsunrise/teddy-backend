package models

import (
	"github.com/mongodb/mongo-go-driver/bson/objectid"
	"time"
)

type Value struct {
	ID    string    `bson:"id"`
	Time  time.Time `bson:"time"`
	Value string    `bson:"value"`
}

type Segment struct {
	ID         objectid.ObjectID `bson:"_id"`
	InfoID     objectid.ObjectID `bson:"infoID"`
	Title      string            `bson:"title"`
	No         uint64            `bson:"no"`
	Labels     []string          `bson:"labels"`
	Values     []Value           `bson:"values"`
	Count      uint64            `bson:"count"`
	WatchCount uint64            `bson:"watchCount"`
}
