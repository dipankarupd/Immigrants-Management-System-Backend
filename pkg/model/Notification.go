package model

import "gopkg.in/mgo.v2/bson"

type Notification struct {
	ID       bson.ObjectId `json:"_id" bson:"_id,omitempty"`
	Receiver string        `json:"receiver" bson:"receiver"`
	Message  string        `json:"message" bson:"message"`
}
