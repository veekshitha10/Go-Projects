package models

import (
	"encoding/json"
	"errors"
)

type Order struct {
	ID          uint    `gorm:"primaryKey" json:"id"`
	UserId      uint    `json:"user_id"`
	SchemeCode  string  `json:"scheme_code"`
	Side        string  `json:"side"`
	Amount      float64 `json:"amount"`
	Status      string  `json:"status"`
	NavUsed     float64 `json:"nav_used"`
	Units       float64 `json:"units"`
	ContractURL string  `json:"contract_url"`
	PlacedAt    int64   `json:"placedat" gorm:"index"`
	ConfirmedAt int64   `json:"confirmedat" gorm:"index"`
}

func (o *Order) Validate() error {
	if o.UserId <= 0 {
		return errors.New("invalid userID")
	}

	if o.SchemeCode == "" {
		return errors.New("invalid scheme code ")

	}
	if o.Side == "" || (o.Side != "buy" && o.Side != "sell") {
		return errors.New("invalid side: must be 'buy' or 'sell'")
	}

	if o.Amount == 0 {
		return errors.New("invalid amount ")

	}
	if o.Units == 0 {
		return errors.New("invalid units ")

	}
	return nil
}
func (u *Order) OrderToBytes() ([]byte, error) {
	bytes, err := json.Marshal(u)
	if err != nil {
		return nil, err
	}
	return bytes, nil
}

type PaymentTable struct {
	Id      uint    `json:"id" gorm:"primaryKey;autoIncrement"`
	OrderId uint    `json:"orderid"`
	Amt     float64 `json:"amt"`
	Status  string  `json:"status"`
}
