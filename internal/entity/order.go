package entity

import "time"

// Order /* */
type Order struct {
	Id        int       `json:"order-id"`
	AccountId int       `json:"account-id"`
	ServiceId int       `json:"service-id"`
	Amount    float64   `json:"money-amount"`
	Status    string    `json:"status"`
	Timecode  time.Time `json:"timecode"`
}
