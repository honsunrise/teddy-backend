package models

import "time"

type KeyValuePair struct {
	Key        string    `json:"key" bson:"key"`
	Value      string    `json:"value" bson:"value"`
	ExpireTime time.Time `json:"expire_time" bson:"expire_time"`
}
