package database

import (
	"errors"
	"trainingmod/models"

	"gorm.io/gorm"
)
type IUserDB interface {
	Create(user *models.UserTable) (*models.UserTable, error)
	GetBy(id uint) (*models.UserTable, error)
		GetOrderBy(id uint) (*models.OrderTable, error)

	 CreateOrder(order *models.OrderTable)(*models.OrderTable,error)
}
type UserDb struct {
	DB *gorm.DB
}

func NewUserDB(db *gorm.DB) IUserDB {
	return &UserDb{db}
}

func (udb *UserDb) Create(user *models.UserTable) (*models.UserTable, error) {
	
	tx := udb.DB.Create(user)
	if tx.Error != nil {
		return nil, tx.Error
	}
	return user, nil
}

func (udb *UserDb) GetBy(id uint) (*models.UserTable, error) {
	user := new(models.UserTable)
	tx := udb.DB.Preload("Orders").First(user, id)
	if tx.Error != nil {
		return nil, tx.Error
	}
	return user, nil
}
func (udb *UserDb) 		GetOrderBy(id uint) (*models.OrderTable, error) {
	order := new(models.OrderTable)
	tx := udb.DB.First(order, id)
	if tx.Error != nil {
		return nil, tx.Error
	}
	return order, nil
}
func (udb *UserDb) GetByLimit(limit, offset int) ([]models.UserTable, error) {
	var users []models.UserTable
	tx := udb.DB.Limit(limit).Offset(offset).Find(&users)
	if tx.Error != nil {
		return nil, tx.Error
	}
	return users, nil
}


func (udb *UserDb)  CreateOrder(order *models.OrderTable)(*models.OrderTable,error) {
	_,err:= udb.GetBy(order.UserId)
	if err!=nil{
		return nil,errors.New("invalid userid")
	}
	tx := udb.DB.Create(order)
	if tx.Error != nil {
		return nil, tx.Error
	}
	return order, nil
}
