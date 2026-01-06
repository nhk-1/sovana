package model

import "time"

// Transaction represents a single bank transaction entry.
type Transaction struct {
	Date   time.Time
	Label  string
	Amount float64
}

// Subscription captures detected recurring payments for a merchant.
type Subscription struct {
	Merchant   string    `json:"merchant"`
	Amount     float64   `json:"amount"`
	Period     string    `json:"period"`
	FirstSeen  time.Time `json:"first_seen"`
	LastSeen   time.Time `json:"last_seen"`
	AnnualCost float64   `json:"annual_cost"`
	CancelURL  string    `json:"cancel_url"`
}
