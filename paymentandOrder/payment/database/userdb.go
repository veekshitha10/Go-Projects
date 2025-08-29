package database

import (
	"trainingmod/models"

	"gorm.io/gorm"
)
type IPaymentDB interface {
	Create(Payment *models.PaymentTable) (*models.PaymentTable, error)
	
}
type PaymentDb struct {
	DB *gorm.DB
}

func NewPaymentDB(db *gorm.DB) IPaymentDB {
	return &PaymentDb{db}
}

func (udb *PaymentDb) Create(Payment *models.PaymentTable) (*models.PaymentTable, error) {
	
	tx := udb.DB.Create(Payment)
	if tx.Error != nil {
		return nil, tx.Error
	}
	return Payment, nil
}
