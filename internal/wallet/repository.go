package wallet

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type WalletRepository interface{
	CreateWallet(ctx context.Context, wallet *Wallet) error
	GetWalletByUserID(ctx context.Context, userID string) ([]*Wallet, error)
	AddBalance(ctx context.Context, userID string, walletType string, balance float64) error
	CreateTransaction(ctx context.Context, transaction *Transactions) error
	GetBalanceUser(ctx context.Context, userID string, walletType string) (*Wallet, error)
	DeductBalance(ctx context.Context, price float64, walletType string, userID string) error
}

type walletRepository struct{
	collection *mongo.Collection
	collectionTransaction *mongo.Collection	
}

func NewWalletRepository(collection *mongo.Collection, collectionTransaction *mongo.Collection) WalletRepository {
	return &walletRepository{
		collection: collection,
		collectionTransaction: collectionTransaction,
	}
}

func (r *walletRepository) CreateWallet(ctx context.Context, wallet *Wallet) error {
	_, err := r.collection.InsertOne(ctx, wallet)
	return err
}

func (r *walletRepository) GetWalletByUserID(ctx context.Context, userID string) ([]*Wallet, error) {

	var wallets []*Wallet

	cursor, err := r.collection.Find(ctx, bson.M{"user_id": userID})
	if err != nil {
		return nil, err
	}

	if err = cursor.All(ctx, &wallets); err != nil {
		return nil, err
	}

	return wallets, nil

}

func (r *walletRepository) AddBalance(ctx context.Context, userID string, walletType string, balance float64) error {
	
	filter := bson.M{
		"user_id": userID,
		"wallet_type": walletType,
	}

	update := bson.M{
		"$set": bson.M{
			"balance": balance,
		},
	}

	_, err := r.collection.UpdateOne(ctx, filter, update)

	return err

}

func (r *walletRepository) CreateTransaction(ctx context.Context, transaction *Transactions) error {
	_, err := r.collectionTransaction.InsertOne(ctx, transaction)
	return err
}

func (r *walletRepository) GetBalanceUser(ctx context.Context, userID string, walletType string) (*Wallet, error) {
	
	var result *Wallet
	
	filter := bson.M{
		"user_id": userID,
		"wallet_type": walletType,
	}

	err := r.collection.FindOne(ctx, filter).Decode(&result)

	return result, err
}

func (r *walletRepository) DeductBalance(ctx context.Context, price float64, walletType string, userID string) error {
	
	filter := bson.M{
		"user_id": userID,
		"wallet_type": walletType,
	}

	update := bson.M{
		"$inc": bson.M{
			"balance": -price,
		},
	}

	_, err := r.collection.UpdateOne(ctx, filter, update)
	return err
}