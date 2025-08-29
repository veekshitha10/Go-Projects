package database

import (
	"trainingmod/models"

	"gorm.io/gorm"
)

type IOrderDB interface {
	Create(user *models.OrdersModel) (*models.OrdersModel, error)
	FetchValues(scrip string) ([]models.OrdersModel, error)
}
type OrderDb struct {
	DB *gorm.DB
}

func NewOrderDB(db *gorm.DB) IOrderDB {
	return &OrderDb{db}
}

func (udb *OrderDb) Create(user *models.OrdersModel) (*models.OrdersModel, error) {

	tx := udb.DB.Create(user)
	if tx.Error != nil {
		return nil, tx.Error
	}
	return user, nil
}

func (udb *OrderDb) FetchValues(scrip string) ([]models.OrdersModel, error) {
	var Orders []models.OrdersModel

	tx := udb.DB.Where("scrip=?", scrip).Find(&Orders)
	if tx.Error != nil {
		return nil, tx.Error
	}
	return Orders, nil
}
