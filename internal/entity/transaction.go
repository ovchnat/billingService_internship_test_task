package entity

import "time"

// Transaction /* */
type Transaction struct {
	Id            int       `json:"tx-id"`
	AccountIdTo   int       `json:"receive-account"`
	AccountIdFrom int       `json:"transfer-account"`
	Amount        float64   `json:"money-amount"`
	Status        string    `json:"status"`
	EventType     string    `json:"event-type"`
	Timecode      time.Time `json:"timecode"`
}
