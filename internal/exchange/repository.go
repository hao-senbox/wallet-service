package exchange

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type ExchangeRepository interface {
	CreateExchangeRate(ctx context.Context, exchangeRate *ExchangeRate) error
	GetAllExchangeRate(ctx context.Context) ([]*ExchangeRate, error)
	GetExchangeRate(ctx context.Context, id primitive.ObjectID) (*ExchangeRate, error)
	GetFirstExchangeRate(ctx context.Context) (*ExchangeRate, error)
	UpdateExchangeRate(ctx context.Context, id primitive.ObjectID, exchangeRate bson.M) error
	DeleteExchangeRate(ctx context.Context, id primitive.ObjectID) error
}

type exchangeRepository struct {
	collection *mongo.Collection	
}

func NewExchangeRepository(collection *mongo.Collection) ExchangeRepository {
	return &exchangeRepository{
		collection: collection,
	}
}

func (r *exchangeRepository) CreateExchangeRate(ctx context.Context, exchangeRate *ExchangeRate) error {
	_, err := r.collection.InsertOne(ctx, exchangeRate)
	return err
}

func (r *exchangeRepository) GetAllExchangeRate(ctx context.Context) ([]*ExchangeRate, error) {

	var exchangeRates []*ExchangeRate

	cursor, err := r.collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}

	if err = cursor.All(ctx, &exchangeRates); err != nil {
		return nil, err
	}

	return exchangeRates, nil
	
}

func (r *exchangeRepository) GetExchangeRate(ctx context.Context, id primitive.ObjectID) (*ExchangeRate, error) {
	var exchangeRate ExchangeRate
	err := r.collection.FindOne(ctx, bson.M{"_id": id}).Decode(&exchangeRate)
	return &exchangeRate, err
}

func (r *exchangeRepository) GetFirstExchangeRate(ctx context.Context) (*ExchangeRate, error) {

	var exchangeRate ExchangeRate

	otps := options.FindOne().SetSort(bson.D{{Key: "created_at", Value: -1}})

	err := r.collection.FindOne(ctx, bson.M{"active": true}, otps).Decode(&exchangeRate)
	if err != nil {
		return nil, err
	}

	return &exchangeRate, err
}

func (r *exchangeRepository) UpdateExchangeRate(ctx context.Context, id primitive.ObjectID, exchangeRate bson.M) error {
	_, err := r.collection.UpdateOne(ctx, bson.M{"_id": id}, bson.M{"$set": exchangeRate})
	return err
}

func (r *exchangeRepository) DeleteExchangeRate(ctx context.Context, id primitive.ObjectID) error {
	_, err := r.collection.DeleteOne(ctx, bson.M{"_id": id})
	return err
}
