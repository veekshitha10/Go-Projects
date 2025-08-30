package database

import (
	//"encoding/json"
	"errors"
	//"fmt"
	"mutualfundminiproject/models"
	// "net/http"
	// "net/url"
	// "strings"

	"gorm.io/gorm"
)

type IUserDB interface {
	Create(user *models.UserTable) (*models.UserTable, error)
	GetOrder(id string) (*models.SchemeTable, error)
	CreateOrder(order *models.Order) (*models.Order, error)
	GetOrdersByUser(userID uint) ([]models.Order, error)
	GetBy(id uint) (*models.UserTable, error)
	UpdateUser(user *models.UserTable) (*models.UserTable, error)
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

	// keycloakURL := "http://keycloak:8080/realms/mutualfund/protocol/openid-connect/token"

	// 	data := url.Values{}
	// 	data.Set("grant_type", "password")
	// 	data.Set("client_id", "frontend") // your Keycloak client ID
	// 	data.Set("username", user.Name)
	// 	data.Set("password", user.Password)

	// 	req, err := http.NewRequest("POST", keycloakURL, strings.NewReader(data.Encode()))
	// 	if err != nil {
	// 		return "", err
	// 	}

	// 	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	// 	client := &http.Client{}
	// 	resp, err := client.Do(req)
	// 	if err != nil {
	// 		return "", err
	// 	}
	// 	defer resp.Body.Close()

	// 	if resp.StatusCode != 200 {
	// 		return "", fmt.Errorf("login failed with status %d", resp.StatusCode)
	// 	}

	// 	var respData struct {
	// 		AccessToken string `json:"access_token"`
	// 	}

	// 	if err := json.NewDecoder(resp.Body).Decode(&respData); err != nil {
	// 		return "", err
	// 	}

	// 	return respData.AccessToken, nil
}

func (udb *UserDb) GetBy(id uint) (*models.UserTable, error) {
	user := new(models.UserTable)
	tx := udb.DB.Preload("Orders").First(user, id)
	if tx.Error != nil {
		return nil, tx.Error
	}
	return user, nil
}
func (udb *UserDb) UpdateUser(user *models.UserTable) (*models.UserTable, error) {
tx:=	udb.DB.Model(&models.UserTable{}).Where("id =?", user.Id).Updates(user)
		if tx.Error != nil {
		return nil, tx.Error
	}
	return user, nil
}
func (udb *UserDb) GetOrderBy(id uint) (*models.Order, error) {
	order := new(models.Order)
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

func (udb *UserDb) CreateOrder(order *models.Order) (*models.Order, error) {
	_, err := udb.GetBy(order.UserId)
	if err != nil {
		return nil, errors.New("invalid userid")
	}
	tx := udb.DB.Create(order)
	if tx.Error != nil {
		return nil, tx.Error
	}
	return order, nil
}
func (udb *UserDb) GetOrder(id string) (*models.SchemeTable, error) {
	scheme := new(models.SchemeTable)
	tx := udb.DB.First(scheme, "scheme_code = ?", id)
	if tx.Error != nil {
		return nil, tx.Error
	}
	return scheme, nil
}
func (udb *UserDb) GetOrdersByUser(userID uint) ([]models.Order, error) {
	var orders []models.Order
	result := udb.DB.Where("user_id = ?", userID).Find(&orders)
	return orders, result.Error
}
