package model

type Person struct {
	ID        interface{} `json:"_id,omitempty" bson:"_id,omitempty"`
	Firstname string      `json:"firstname,omitempty" bson:"firstname,omitempty"`
	Lastname  string      `json:"lastname,omitempty" bson:"lastname,omitempty"`
}

type Persons []Person