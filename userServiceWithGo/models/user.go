package models

import (
	"encoding/json"
	"errors"
)

var (
	ErrInvalidName   = errors.New("invalid name field")
	ErrInvalidEmail  = errors.New("invalid email field")
	ErrInvalidMobile = errors.New("invalid mobile field")
		Job   chan uint

)

type UserTable struct {
	CommonModel        // promoted field
	Name        string `json:"name"`
	Email       string `json:"email" gorm:"uniqueIndex"`
	Orders      []OrderTable `json:"orders ,omitempty" gorm:"foreignKey:UserId"`
}

func (u *UserTable) Validate() error {
	if u.Name == "" {
		return ErrInvalidName
	}
	if u.Email == "" {
		return ErrInvalidEmail
	}
	
	return nil
}
func (u *UserTable) ToBytes() []byte {
	bytes, _ := json.Marshal(u)
	return bytes
}


