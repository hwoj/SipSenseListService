package models

import "time"

type List struct {
	ID            string    `json:"id,omitempty" bson:"_id,omitempty"`
	Name          string    `json:"name,omitempty" bson:"name,omitempty"`
	Owner         string    `json:"owner,omitempty" bson:"owner,omitempty"`
	Public        bool      `json:"public" bson:"pubpic"`
	Collaborators []string  `json:"collaborators,omitempty" bson:"collaborators,omitempty"`
	Drinks        []string  `json:"drinks" bson:"drinks"`
	LastModified  time.Time `json:"lastModified" bson:"lastModified"`
}
