package models

import (
	"errors"
	"strings"
)

type OrdersModel struct {
	CommonModel
	Scrip       string  `json:"scrip"`
	Quantity    int     `json:"quantity"`
	Price       float64 `json:"price"`
	BuysellType string  `json:"buyselltype"`
}

var (
	ErrInvalidScrip       = errors.New("invalid scrip name")
	ErrInvalidQuantity    = errors.New("invalid email")
	ErrInvalidPrice       = errors.New("invalid mobile ")
	ErrInvalidBuySellType = errors.New("invalid BuysellType ")
)

func (u *OrdersModel) Validate() error {
	if u.Scrip == "" {
		return ErrInvalidScrip
	}
	if u.BuysellType == "" || (strings.ToLower(u.BuysellType) != "buy" && strings.ToLower(u.BuysellType) != "sell") {
		return ErrInvalidBuySellType
	}
	if u.Price == 0 {
		return ErrInvalidPrice
	}
	if u.Quantity == 0 {
		return ErrInvalidQuantity
	}
	return nil
}
