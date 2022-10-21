package entity

import "time"

// Account /* */
type Account struct {
	Id            int       `json:"order-id"`
	UserId        int       `json:"user-id" binding:"required"`
	CurrAmount    float64   `json:"money-amount" binding:"required"`
	PendingAmount string    `json:"status"`
	LastUpdatedAt time.Time `json:"last-updated-time"`
}

type GetBalanceRequest struct {
	UserId int `json:"user-id" binding:"required"`
}

type GetBalanceResponse struct {
	Balance int64 `json:"user-balance"`
	Pending int64 `json:"user-pending-amount"`
}

type UpdateBalanceRequest struct {
	UserId int   `json:"user-id" binding:"required"`
	Sum    int64 `json:"update-amount" binding:"required"`
}
type UpdateBalanceResponse struct {
	AccountId int   `json:"account-id" binding:"required"`
	UserId    int   `json:"user-id" binding:"required"`
	Balance   int64 `json:"current-balance" binding:"required"`
	Pending   int64 `json:"user-pending-amount"`
}
