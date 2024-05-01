package app

import "go.mongodb.org/mongo-driver/bson/primitive"

type Product struct {
	ObjectId    primitive.ObjectID `json: "id" bson: "_id"`
	Name        string             `json: "name" bson: "name,omitempty"`
	Description string             `json: "description" bson: "description,omitempty"`
	Category    string             `json: "category" bson: "category,omitempty"`
	Price       int                `json: "price" bson: "price,omitempty"`
}

type Products = []Product
