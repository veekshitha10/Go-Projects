package models

import (
	"encoding/json"
	"errors"
)

var (
	ErrInvalidName   = errors.New("invalid name field")
	ErrInvalidEmail  = errors.New("invalid email field")
	ErrInvalidMobile = errors.New("invalid mobile field")
	Job              chan uint
)

type PaymentTable struct {
	Id      uint    `json:"id" gorm:"primaryKey;autoIncrement"`
	OrderId uint    `json:"orderid"`
	Amt     float64 `json:"amt"`
	Status  string  `json:"status"`
}

func (u *PaymentTable) ToBytes() []byte {
	bytes, _ := json.Marshal(u)
	return bytes
}


type OrderTable struct {
	CommonModel
    Status string `json:"status"`
	UserId   uint   `json:"user_id"`
	Totalcents int `json:"total_cents"`
	Orderstatus string `json:"order_status"`
}
