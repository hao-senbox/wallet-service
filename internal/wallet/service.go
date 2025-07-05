package wallet

import (
	"context"
	"fmt"
	"time"
	"wallet-service/internal/exchange"
	"wallet-service/internal/user"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type WalletService interface {
	CreateWallet(ctx context.Context, userID string) error
	GetWalletByUserID(ctx context.Context, userID string) (*WalletByUser, error)
	AddBalance(ctx context.Context, req *AddBalanceRequest, userID string) error
	DeductBalance(ctx context.Context, req *DeductBalanceRequest, userID string) error // Fixed typo
}

type walletService struct {
	walletRepo      WalletRepository
	exchangeService exchange.ExchangeService
	userService     user.UserService
}

func NewWalletService(walletRepo WalletRepository, exchangeService exchange.ExchangeService, userService user.UserService) WalletService {
	return &walletService{
		walletRepo:      walletRepo,
		exchangeService: exchangeService,
		userService:     userService,
	}
}

func (s *walletService) CreateWallet(ctx context.Context, userID string) error {
	// Validate input
	if userID == "" {
		return fmt.Errorf("user_id is required")
	}

	// Check if wallets already exist
	existingWallets, err := s.walletRepo.GetWalletByUserID(ctx, userID)
	if err == nil && len(existingWallets) > 0 {
		return fmt.Errorf("wallet already exists for user %s", userID)
	}

	// Create store wallet
	walletStore := &Wallet{
		ID:         primitive.NewObjectID(),
		UserID:     userID,
		WalletType: "store",
		Balance:    0,
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}

	err = s.walletRepo.CreateWallet(ctx, walletStore)
	if err != nil {
		return fmt.Errorf("failed to create store wallet: %w", err)
	}

	// Create service wallet
	walletService := &Wallet{
		ID:         primitive.NewObjectID(),
		UserID:     userID,
		WalletType: "service",
		Balance:    0,
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}

	err = s.walletRepo.CreateWallet(ctx, walletService)
	if err != nil {
		return fmt.Errorf("failed to create service wallet: %w", err)
	}

	return nil
}

func (s *walletService) GetWalletByUserID(ctx context.Context, userID string) (*WalletByUser, error) {
	// Validate input
	if userID == "" {
		return nil, fmt.Errorf("user_id is required")
	}
	
	wallet, err := s.walletRepo.GetWalletByUserID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get wallet: %w", err)
	}

	// Create wallet if doesn't exist
	if len(wallet) == 0 {
		err := s.CreateWallet(ctx, userID)
		if err != nil {
			return nil, fmt.Errorf("failed to create wallet: %w", err)
		}

		// Get wallet again after creation
		wallet, err = s.walletRepo.GetWalletByUserID(ctx, userID)
		if err != nil {
			return nil, fmt.Errorf("failed to get wallet after creation: %w", err)
		}
	}

	var walletUsers []WalletUser
	for _, wal := range wallet {
		walletUsers = append(walletUsers, WalletUser{
			Balance:    wal.Balance,
			WalletType: wal.WalletType,
		})
	}

	return &WalletByUser{
		UserID: userID,
		Wallet: walletUsers,
	}, nil
}

func (s *walletService) AddBalance(ctx context.Context, req *AddBalanceRequest, userID string) error {
	// Validate input
	if req == nil {
		return fmt.Errorf("request is required")
	}

	if req.Balance <= 0 {
		return fmt.Errorf("balance must be greater than 0")
	}

	if req.UserID == "" {
		return fmt.Errorf("user_id is required")
	}

	if req.WalletType == "" {
		return fmt.Errorf("wallet_type is required")
	}

	if userID == "" {
		return fmt.Errorf("admin user_id is required")
	}

	// Validate wallet type
	if req.WalletType != "store" && req.WalletType != "service" {
		return fmt.Errorf("invalid wallet_type: %s", req.WalletType)
	}

	// Check if user's wallet exists
	_, err := s.GetWalletByUserID(ctx, req.UserID)
	if err != nil {
		return fmt.Errorf("failed to get user wallet: %w", err)
	}

	// Get exchange rate
	exchange, err := s.exchangeService.GetFirstExchangeRate(ctx)
	if err != nil {
		return fmt.Errorf("failed to get exchange rate: %w", err)
	}

	if exchange.Rate <= 0 {
		return fmt.Errorf("invalid exchange rate: %f", exchange.Rate)
	}

	// Get current balance
	wallet, err := s.walletRepo.GetBalanceUser(ctx, req.UserID, req.WalletType)
	if err != nil {
		return fmt.Errorf("failed to get current balance: %w", err)
	}

	// Calculate new balance (check for overflow)
	addAmount := exchange.Rate * req.Balance
	newBalance := wallet.Balance + addAmount

	// Basic overflow check
	if newBalance < wallet.Balance {
		return fmt.Errorf("balance overflow detected")
	}

	// Update balance
	err = s.walletRepo.AddBalance(ctx, req.UserID, req.WalletType, newBalance)
	if err != nil {
		return fmt.Errorf("failed to update balance: %w", err)
	}

	// Create transaction record
	transaction := &Transactions{
		ID:        primitive.NewObjectID(),
		UserID:    req.UserID,
		Type:      "deposit",
		Money:     &req.Balance,
		Amount:    addAmount,
		Currency:  &exchange.Currency,
		AdminID:   &userID,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(), // Added UpdatedAt
	}

	err = s.walletRepo.CreateTransaction(ctx, transaction)
	if err != nil {
		// TODO: Consider rollback balance update
		return fmt.Errorf("failed to create transaction: %w", err)
	}

	return nil
}

func (s *walletService) DeductBalance(ctx context.Context, req *DeductBalanceRequest, userID string) error {
	// Validate input
	if req == nil {
		return fmt.Errorf("request is required")
	}

	if userID == "" {
		return fmt.Errorf("user_id is required")
	}

	if req.PriceService < 0 || req.PriceStore < 0 {
		return fmt.Errorf("prices cannot be negative")
	}

	if req.PriceService == 0 && req.PriceStore == 0 {
		return fmt.Errorf("at least one price must be greater than 0")
	}

	// Get exchange rate
	exchange, err := s.exchangeService.GetFirstExchangeRate(ctx)
	if err != nil {
		return fmt.Errorf("failed to get exchange rate: %w", err)
	}

	_, err = s.GetWalletByUserID(ctx, userID)
	if err != nil {
		return err
	}

	dataUser, err := s.userService.GetAllUser(ctx)
	if err != nil {
		return err
	}

	found := false
	for _, item := range dataUser {
		if item.UserID == userID {
			found = true
			break
		}
	}

	if !found {
		return fmt.Errorf("user id not found")
	}

	// Check service wallet balance if needed
	if req.PriceService > 0 {
		serviceWallet, err := s.walletRepo.GetBalanceUser(ctx, userID, "service")
		if err != nil {
			return fmt.Errorf("failed to get service wallet balance: %w", err)
		}

		if serviceWallet.Balance < req.PriceService {
			return fmt.Errorf("insufficient service wallet balance: have %f, need %f",
				serviceWallet.Balance, req.PriceService)
		}
	}

	// Check store wallet balance if needed
	if req.PriceStore > 0 {
		storeWallet, err := s.walletRepo.GetBalanceUser(ctx, userID, "store")
		if err != nil {
			return fmt.Errorf("failed to get store wallet balance: %w", err)
		}

		if storeWallet.Balance < req.PriceStore {
			return fmt.Errorf("insufficient store wallet balance: have %f, need %f",
				storeWallet.Balance, req.PriceStore)
		}
	}

	// Deduct from service wallet
	if req.PriceService > 0 {
		err = s.walletRepo.DeductBalance(ctx, req.PriceService, "service", userID)
		if err != nil {
			return fmt.Errorf("failed to deduct from service wallet: %w", err)
		}
	}

	// Deduct from store wallet
	if req.PriceStore > 0 {
		err = s.walletRepo.DeductBalance(ctx, req.PriceStore, "store", userID)
		if err != nil {
			// TODO: Rollback service wallet deduction
			return fmt.Errorf("failed to deduct from store wallet: %w", err)
		}
	}

	// Create transaction record
	totalAmount := req.PriceService + req.PriceStore
	transaction := &Transactions{
		ID:        primitive.NewObjectID(),
		UserID:    userID,
		Type:      "purchase",
		Amount:    totalAmount,
		Money:     nil,
		Currency:  &exchange.Currency,
		AdminID:   nil,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	err = s.walletRepo.CreateTransaction(ctx, transaction)
	if err != nil {
		return fmt.Errorf("failed to create transaction: %w", err)
	}

	return nil
}
