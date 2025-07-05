package exchange

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ExchangeRate struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	Currency  string             `bson:"currency" json:"currency"`
	Rate      float64            `bson:"rate" json:"rate"`
	Active    bool               `bson:"active" json:"active"`
	CreatedAt time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt time.Time          `bson:"updated_at" json:"updated_at"`
}

