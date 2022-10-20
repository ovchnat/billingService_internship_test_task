package repository

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

type BillingRepo struct {
	DepositMoney
	WithdrawMoney
	GetBalance
	ReserveAmount
	OrderConfirm
	Transfer
}

func NewRepo() *BillingRepo {
	return &BillingRepo{}
}
