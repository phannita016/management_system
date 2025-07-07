package store

import (
	"context"
	"errors"
	"time"

	"github.com/phannita016/management/dtos"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type (
	Management interface {
		Create(ctx context.Context, body dtos.ManagementRequest) error
		Update(ctx context.Context, body dtos.Management) error
		Delete(ctx context.Context, filter bson.M) error
		GetAll(ctx context.Context) ([]dtos.Management, error)
		FindByID(ctx context.Context, id primitive.ObjectID) (dtos.Management, error)
	}

	managementEntity struct {
		collection *mongo.Collection
	}
)

func NewDatabaseManagement(ctx context.Context, database *mongo.Client) Management {
	return &managementEntity{
		collection: database.Database("management").Collection("managements"),
	}
}

func (m *managementEntity) Create(ctx context.Context, body dtos.ManagementRequest) error {
	model := dtos.Management{
		ID:         primitive.NewObjectID(),
		Name:       body.Name,
		Nickname:   body.Nickname,
		Gender:     body.Gender,
		Age:        body.Age,
		Role:       body.Role,
		CreateDate: time.Now(),
	}

	_, err := m.collection.InsertOne(ctx, model)
	return err
}

func (m *managementEntity) Update(ctx context.Context, model dtos.Management) error {
	filter := bson.M{"_id": model.ID}
	update := bson.M{
		"$set": bson.M{
			"name":       model.Name,
			"nickname":   model.Nickname,
			"gender":     model.Gender,
			"age":        model.Age,
			"role":       model.Role,
			"updateDate": time.Now(),
		},
	}
	_, err := m.collection.UpdateOne(ctx, filter, update)
	return err
}

func (m *managementEntity) Delete(ctx context.Context, filter bson.M) error { // change here
	_, err := m.collection.DeleteOne(ctx, filter)
	return err
}

func (m *managementEntity) GetAll(ctx context.Context) ([]dtos.Management, error) {
	opts := options.Find().SetSort(bson.D{{Key: "createDate", Value: -1}})
	cur, err := m.collection.Find(ctx, bson.M{}, opts)
	if err != nil {
		return nil, err
	}
	defer cur.Close(ctx)

	var results []dtos.Management
	for cur.Next(ctx) {
		var model dtos.Management
		if err = cur.Decode(&model); err != nil {
			return nil, err
		}
		results = append(results, model)
	}

	return results, nil
}

func (m *managementEntity) FindByID(ctx context.Context, id primitive.ObjectID) (dtos.Management, error) {
	var model dtos.Management
	err := m.collection.FindOne(ctx, bson.M{"_id": id}).Decode(&model)
	if errors.Is(err, mongo.ErrNoDocuments) {
		return model, mongo.ErrNoDocuments
	}
	return model, err
}
