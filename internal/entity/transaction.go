package entity

import "time"

// Transaction /* */
type Transaction struct {
	Id            int64     `json:"tx-id"`
	AccountIdTo   int64     `json:"receive-account"`
	AccountIdFrom int64     `json:"transfer-account"`
	Amount        int64     `json:"money-amount"`
	Status        string    `json:"status"`
	Timecode      time.Time `json:"timecode"`
}

type TransferRequest struct {
	SenderId   int64 `json:"sender-id" binding:"required"`
	ReceiverId int64 `json:"receiver-id" binding:"required"`
	Sum        int64 `json:"transfer-amount" binding:"required"`
}

type TransferResponse struct {
	AccountIdTo   int64     `json:"receive-account"`
	AccountIdFrom int64     `json:"transfer-account"`
	Amount        int64     `json:"money-amount"`
	Status        string    `json:"status"`
	EventType     string    `json:"event-type"`
	Timecode      time.Time `json:"created-at"`
}
