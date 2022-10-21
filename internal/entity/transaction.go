package entity

import "time"

// Transaction /* */
type Transaction struct {
	Id            int       `json:"tx-id"`
	AccountIdTo   int       `json:"receive-account"`
	AccountIdFrom int       `json:"transfer-account"`
	Amount        float64   `json:"money-amount"`
	Status        string    `json:"status"`
	Timecode      time.Time `json:"timecode"`
}

type TransferRequest struct {
	SenderId   int   `json:"sender-id" binding:"required"`
	ReceiverId int   `json:"receiver-id" binding:"required"`
	Sum        int64 `json:"transfer-amount" binding:"required"`
}

type TransferResponse struct {
	Message     string `json:"message"`
	Transaction `json:"transaction-info"`
}
