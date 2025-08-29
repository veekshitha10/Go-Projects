package models

import (
	"encoding/json"
	"errors"
)

type OrderTable struct {
	CommonModel
    Status string `json:"status"`
	UserId   uint   `json:"user_id"`
	Totalcents int `json:"total_cents"`
	Orderstatus string `json:"order_status"`
}

func (o *OrderTable) Validate() error {
	if o.UserId <= 0 {
		return errors.New("invalid userID")
	} 
	if o.Totalcents ==0 {
		return errors.New("invalid total cent value")
	} 
	return nil //interface can be nil
}
func (u *OrderTable) OrderToBytes() ([]byte,error) {
	bytes, err := json.Marshal(u)
	if err !=nil{
		return  nil,err
	}
	return bytes,nil
}

type PaymentTable struct {
	Id      uint    `json:"id" gorm:"primaryKey;autoIncrement"`
	OrderId uint    `json:"orderid"`
	Amt     float64 `json:"amt"`
	Status  string  `json:"status"`
}
