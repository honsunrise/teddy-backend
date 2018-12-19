package models

import (
	"github.com/mongodb/mongo-go-driver/bson/objectid"
	"time"
)

type BehaviorInfoItem struct {
	InfoId objectid.ObjectID `bson:"infoId"`
	Time   time.Time         `bson:"time"`
}

type BehaviorUserItem struct {
	UID  string    `bson:"uid"`
	Time time.Time `bson:"time"`
}

type Behavior struct {
	Id        objectid.ObjectID  `bson:"_id"`
	UID       string             `bson:"uid"`
	FirstTime time.Time          `bson:"firstTime"`
	LastTime  time.Time          `bson:"lastTime"`
	Count     uint64             `bson:"count"`
	Items     []BehaviorInfoItem `bson:"items"`
}
