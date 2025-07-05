package exchange

type CreateExchangeRateRequest struct {
	Currency string  `json:"currency"`
	Rate     float64 `json:"rate"`
	Active   bool    `json:"active"`
}

type UpdateExchangeRateRequest struct {
	Currency *string  `json:"currency"`
	Rate     *float64 `json:"rate"`
	Active   *bool    `json:"active"`
}