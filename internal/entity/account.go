package entity

import "time"

// Account /* */
type Account struct {
	Id            int       `json:"order-id"`
	UserId        int       `json:"user-id"`
	CurrAmount    float64   `json:"money-amount"`
	PendingAmount string    `json:"status"`
	LastUpdatedAt time.Time `json:"last-updated-time"`
}
