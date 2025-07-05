package wallet

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Wallet struct {
	ID         primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	UserID     string             `bson:"user_id" json:"user_id"`
	Balance    float64            `bson:"balance" json:"balance"`
	WalletType string             `bson:"wallet_type" json:"wallet_type"`
	CreatedAt  time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt  time.Time          `bson:"updated_at" json:"updated_at"`
}

type Transactions struct {
	ID        primitive.ObjectID  `bson:"_id,omitempty" json:"id,omitempty"`
	UserID    string              `bson:"user_id" json:"user_id"`
	Type      string              `bson:"type" json:"type"`
	Amount    float64             `bson:"amount" json:"amount"`
	Money     *float64            `bson:"money" json:"money"`
	Currency  *string             `bson:"currency" json:"currency"`
	OrderID   *primitive.ObjectID `bson:"order_id" json:"order_id"`
	AdminID   *string             `bson:"admin_id" json:"admin_id"`
	CreatedAt time.Time           `bson:"created_at" json:"created_at"`
	UpdatedAt time.Time           `bson:"updated_at" json:"updated_at"`
}

type ServiceUsageLog struct {
	ID             primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	OrganizationID primitive.ObjectID `bson:"organization_id" json:"organization_id"`
	UserID         string             `bson:"user_id" json:"user_id"`
	ServiceName    string             `bson:"service_name" json:"service_name"`
	Action         string             `bson:"action" json:"action"`
	Timestamp      time.Time          `bson:"timestamp" json:"timestamp"`
}
