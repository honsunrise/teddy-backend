package models

import (
	"github.com/mongodb/mongo-go-driver/bson/primitive"
	"time"
)

type BehaviorInfoItem struct {
	InfoId primitive.ObjectID `bson:"infoId"`
	Time   time.Time          `bson:"time"`
}

type BehaviorUserItem struct {
	UID  string    `bson:"uid"`
	Time time.Time `bson:"time"`
}

type Behavior struct {
	Id        primitive.ObjectID  `bson:"_id"`
	UID       string              `bson:"uid"`
	FirstTime time.Time           `bson:"firstTime"`
	LastTime  time.Time           `bson:"lastTime"`
	Count     uint64              `bson:"count"`
	Items     []*BehaviorInfoItem `bson:"items"`
}
