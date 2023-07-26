package model

import (
	"gopkg.in/mgo.v2/bson"
)

type Feedback struct {
	ID          bson.ObjectId `json:"_id" bson:"_id,omitempty"`
	PassportNo  int           `json:"passportno" bson:"passportno"`
	Comment     string        `json:"comment" bson:"comment"`
	Rating      int           `json:"rating" bson:"rating"`
	ImmigrantID bson.ObjectId `json:"immigrant,omitempty" bson:"immigrant,omitempty"`
}
