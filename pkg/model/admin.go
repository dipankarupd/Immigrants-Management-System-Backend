package model

import "gopkg.in/mgo.v2/bson"

type Admin struct {
	ID       bson.ObjectId `json:"_id" bson:"_id,omitempty"`
	Username string        `json:"username" bson:"username"`
}
