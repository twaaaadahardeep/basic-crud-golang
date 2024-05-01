package app

import (
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type DBRepository struct {
	client *mongo.Client
	db     *mongo.Database
	col    *mongo.Collection
}

type Repository interface {
	findAll(ctx context.Context) (*Products, error)
	findOne(ctx context.Context, id string) (*Product, error)
	insertOne(ctx context.Context, product *Product) error
	updateOne(ctx context.Context, id string, product *Product) error
	updateField(ctx context.Context, id, field string, value interface{}) error
	deleteOne(ctx context.Context, id string) error
	Close()
}

func Connect() (Repository, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		log.Fatal(err)
	}

	db := client.Database("Shop")
	col := db.Collection("Products")

	return &DBRepository{
		client: client,
		db:     db,
		col:    col,
	}, nil
}

func (d *DBRepository) findAll(ctx context.Context) (*Products, error) {
	c, err := d.col.Find(ctx, bson.D{})
	if err != nil {
		return nil, err
	}

	var ps Products

	for c.Next(ctx) {
		var p Product
		if err := c.Decode(&p); err != nil {
			return nil, err
		}

		ps = append(ps, p)
	}

	if err := c.Close(ctx); err != nil {
		return nil, err
	}

	return &ps, nil
}

func (d *DBRepository) findOne(ctx context.Context, id string) (*Product, error) {
	obId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	res := d.col.FindOne(ctx, bson.M{"_id": obId})
	var p Product
	if err := res.Decode(&p); err != nil {
		return nil, err
	}

	return &p, nil
}

func (d *DBRepository) insertOne(ctx context.Context, product *Product) error {
	product.ObjectId = primitive.NewObjectID()
	if _, err := d.col.InsertOne(ctx, product); err != nil {
		return err
	}

	return nil
}

func (d *DBRepository) updateOne(ctx context.Context, id string, product *Product) error {
	obId, _ := primitive.ObjectIDFromHex(id)
	if _, err := d.col.ReplaceOne(ctx, bson.M{"_id": obId}, product); err != nil {
		return err
	}

	return nil
}

func (d *DBRepository) updateField(ctx context.Context, id, field string, value interface{}) error {
	obId, _ := primitive.ObjectIDFromHex(id)
	filter := bson.M{"_id": obId}
	update := bson.D{{
		"$set", bson.D{{
			field,
			value,
		}},
	}}

	if _, err := d.col.ReplaceOne(ctx, filter, update); err != nil {
		return err
	}

	return nil
}

func (d *DBRepository) deleteOne(ctx context.Context, id string) error {
	obId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	if _, err := d.col.DeleteOne(ctx, bson.M{"_id": obId}); err != nil {
		return err
	}

	return nil
}

func (d *DBRepository) Close() {
	if err := d.client.Disconnect(context.Background()); err != nil {
		log.Fatal(err)
	}
}
