package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Character struct {
	ID        primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Name      string             `json:"name" validate:"required"`
	Birthday  string             `json:"birthday,omitempty"`
	Dead      bool               `json:"dead" bson:"dead,omitempty"`
	Relevance string             `json:"relevance" validate:"required"`
	Seasons   int                `json:"seasons" validate:"required"`
}
