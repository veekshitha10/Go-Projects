package models

import "errors"

type OrderTable struct {
	CommonModel
    Status string `json:"status"`
	UserId   uint   `json:"user_id"`
	Totalcents int `json:"total_cents"`
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