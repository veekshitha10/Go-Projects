package models

import (
	"encoding/json"
	"errors"
	
)

var (
	ErrInvalidName     = errors.New("invalid name field")
	ErrInvalidPassword = errors.New("invalid password ")
	Channel          chan uint
)

type UserTable struct {
	Id       uint    `json:"id" gorm:"primaryKey;autoIncrement"`
	Name     string  `json:"name"`
	Password string  `json:"password"`
	Orders   []Order `json:"order ,omitempty" gorm:"foreignKey:UserId"`
}

func (u *UserTable) Validate() error {
	if u.Name == "" {
		return ErrInvalidName
	}
	if u.Password == "" {
		return ErrInvalidPassword
	}

	return nil
}
func (u *UserTable) ToBytes() []byte {
	bytes, _ := json.Marshal(u)
	return bytes
}
