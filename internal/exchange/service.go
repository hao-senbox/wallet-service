package exchange

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ExchangeService interface {
	CreateExchangeRate(ctx context.Context, req *CreateExchangeRateRequest) error
	GetAllExchangeRate(ctx context.Context) ([]*ExchangeRate, error)
	GetExchangeRate(ctx context.Context, id string) (*ExchangeRate, error)
	GetFirstExchangeRate(ctx context.Context) (*ExchangeRate, error)
	UpdateExchangeRate(ctx context.Context, id string, req *UpdateExchangeRateRequest) error
	DeleteExchangeRate(ctx context.Context, id string) error
}

type exchangeService struct {
	exchangeRepo ExchangeRepository	
}

func NewExchangeService(exchangeRepo ExchangeRepository) ExchangeService {
	return &exchangeService{
		exchangeRepo: exchangeRepo,
	}
}

func (s *exchangeService) CreateExchangeRate(ctx context.Context,req *CreateExchangeRateRequest) error {
	
	if req.Currency == "" {
		return fmt.Errorf("currency is required")
	}

	if req.Rate == 0 {
		return fmt.Errorf("rate is required")
	}

	exchangeRate := &ExchangeRate {
		ID: primitive.NewObjectID(),
		Currency: req.Currency,
		Rate: req.Rate,
		Active: true,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	return s.exchangeRepo.CreateExchangeRate(ctx, exchangeRate)

}

func (s *exchangeService) GetAllExchangeRate(ctx context.Context) ([]*ExchangeRate, error) {
	return s.exchangeRepo.GetAllExchangeRate(ctx)
}

func (s *exchangeService) GetExchangeRate(ctx context.Context, id string) (*ExchangeRate, error) {

	if id == "" {
		return nil, fmt.Errorf("id is required")
	}

	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	return s.exchangeRepo.GetExchangeRate(ctx, objectID)

}

func (s *exchangeService) GetFirstExchangeRate(ctx context.Context) (*ExchangeRate, error) {
	return s.exchangeRepo.GetFirstExchangeRate(ctx)
}

func (s *exchangeService) UpdateExchangeRate(ctx context.Context, id string, req *UpdateExchangeRateRequest) error {

	if id == "" {
		return fmt.Errorf("id is required")
	}

	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	updateExchangeRate := bson.M{}

	if req.Currency != nil {
		updateExchangeRate["currency"] = req.Currency
	}

	if req.Rate != nil {
		updateExchangeRate["rate"] = req.Rate
	}

	if req.Active != nil {
		updateExchangeRate["active"] = req.Active
	}

	updateExchangeRate["updated_at"] = time.Now()

	if len(updateExchangeRate) == 1 {
		return fmt.Errorf("no fields to update")
	}

	return s.exchangeRepo.UpdateExchangeRate(ctx, objectID, updateExchangeRate)
	
}

func (s *exchangeService) DeleteExchangeRate(ctx context.Context, id string) error {

	if id == "" {
		return fmt.Errorf("id is required")
	}

	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	return s.exchangeRepo.DeleteExchangeRate(ctx, objectID)
	
}