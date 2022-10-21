package entity

import "time"

// Order /* */
type Order struct {
	Id        int       `json:"order-id"`
	UserId    int       `json:"user-id"`
	ServiceId int       `json:"service-id"`
	OrderId   int       `json:"order-id"`
	Amount    float64   `json:"money-amount"`
	Status    string    `json:"status"`
	Timecode  time.Time `json:"timecode"`
}

// ReserveServiceFeeRequest
// Метод резервирования средств с основного баланса на отдельном счете.
// Принимает id пользователя, ИД услуги, ИД заказа, стоимость.

type ReserveServiceFeeRequest struct {
	UserId    int     `json:"user-id"`
	ServiceId int     `json:"service-id"`
	OrderId   int     `json:"order-id"`
	Fee       float64 `json:"fee"`
}

type ReserveServiceFeeResponse struct {
	UserId    int       `json:"user-id"`
	ServiceId int       `json:"service-id"`
	Invoice   float64   `json:"invoice"`
	Status    int       `json:"status"`
	CreatedAt time.Time `json:"created-at"`
}
