package entity

import "time"

// Account /* */
type Account struct {
	Id            int64     `json:"order-id"`
	UserId        int64     `json:"user-id" binding:"required"`
	CurrAmount    float64   `json:"money-amount" binding:"required"`
	PendingAmount string    `json:"pending-amount"`
	LastUpdatedAt time.Time `json:"last-updated-time"`
}

type GetBalanceRequest struct {
	UserId int64 `json:"user-id" binding:"required"`
}

type GetBalanceResponse struct {
	Balance int64 `json:"user-balance"`
	Pending int64 `json:"user-pending-amount"`
}

type UpdateBalanceRequest struct {
	UserId int64 `json:"user-id" binding:"required"`
	Sum    int64 `json:"update-amount" binding:"required"`
}
type UpdateBalanceResponse struct {
	AccountId int64     `json:"account-id" binding:"required"`
	Sum       int64     `json:"sum" binding:"required"`
	Status    string    `json:"operation-status" binding:"required"`
	EventType string    `json:"event-type"`
	CreatedAt time.Time `json:"created-at"`
}
