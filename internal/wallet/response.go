package wallet

type WalletByUser struct {
	UserID string `json:"user_id"`
	Wallet []WalletUser `json:"wallet"`
}

type WalletUser struct {
	Balance    float64            `json:"balance"`
	WalletType string             `json:"wallet_type"`
}
