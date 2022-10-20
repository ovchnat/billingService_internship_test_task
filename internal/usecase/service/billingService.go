package service

import "billingService/internal/usecase/repository"

type DepositMoney interface {
}

type WithdrawMoney interface {
}

type GetBalance interface {
}

type ReserveAmount interface {
}

type OrderConfirm interface {
}

type Transfer interface {
}

type BillingService struct {
	DepositMoney
	WithdrawMoney
	GetBalance
	ReserveAmount
	OrderConfirm
	Transfer
}

func NewService(repos *repository.BillingRepo) *BillingService {
	return &BillingService{}
}
