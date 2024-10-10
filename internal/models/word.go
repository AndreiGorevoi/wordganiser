package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Word struct {
	ID         primitive.ObjectID `bson:"_id,omitempty"`
	Word       string             `bson:"word"`
	Definition string             `bson:"definition"`
	Usage      string             `bson:"usage"`
}
