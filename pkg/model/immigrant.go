package model

import (
	"gopkg.in/mgo.v2/bson"
)

type Immigrant struct {
	ID           bson.ObjectId `json:"id" bson:"_id"`
	Name         string        `json:"name" bson:"name"`
	PassPort_No  int           `json:"passportno" bson:"passportno"`
	Email        string        `json:"email" bson:"email"`
	Gender       string        `json:"gender" bson:"gender"`
	Country      string        `json:"country" bson:"country"`
	Age          int           `json:"age" bson:"age"`
	Arrival_Date string        `json:"arrivaldate" bson:"arrivaldate"`
	Stay_Time    string        `json:"staytime" bson:"staytime"`
	Visa_Type    string        `json:"visatype" bson:"visatype"`
	Approval     *string       `json:"approval,omitempty" bson:"approval,omitempty" validate:"eq=pending | eq=approved | eq=rejected" `
}

func (i *Immigrant) SetDefaultValue() {
	approval := "pending"
	i.Approval = &approval
}
