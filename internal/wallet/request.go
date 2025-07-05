package wallet

type AddBalanceRequest struct {
	UserID     string  `json:"user_id"`
	WalletType string  `json:"wallet_type"`
	Balance    float64 `json:"balance"`
}

type DeductBalanceRequest struct {
	PriceStore   float64 `json:"price_store"`
	PriceService float64 `json:"price_service"`
}
