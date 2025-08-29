package handlers

import (
	"errors"
	"fmt"
	"strings"
	"time"
	"trainingmod/database"
	"trainingmod/models"

	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog/log"
)

type OrderHandler struct {
	database.IOrderDB // prmoted field
}

type IOrderHandler interface {
	CreateOrder(c *fiber.Ctx) error
	CalculateNet(c *fiber.Ctx) error
}

func NewOrderHandler(iOrderdb database.IOrderDB) IOrderHandler {
	return &OrderHandler{iOrderdb}
}

func (uh *OrderHandler) CreateOrder(c *fiber.Ctx) error {
	Order := new(models.OrdersModel)
	err := c.BodyParser(Order)
	if err != nil {
		return err
	}

	err = Order.Validate()
	if err != nil {
		return err
	}
	Order.LastModified = time.Now().Unix()

	Order, err = uh.Create(Order)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "invalid order request")

	}
	fmt.Printf("%s %d %s at ₹%.2f\n",Order.BuysellType,Order.Quantity,Order.Scrip,Order.Price)
	return c.JSON(Order)

}

func (uh *OrderHandler) CalculateNet(c *fiber.Ctx) error {
	scrip := c.Params("scrip")

	if scrip == "" {
		return errors.New("invalid scrip")
	}

	users, err := uh.FetchValues(scrip)
	if err != nil {
		log.Err(err).Msg("data might not be available or some sql issue")
		return errors.New("something went wrong or no data available with that scrip")
	}
	if len(users) <= 0 {
		return errors.New(" no data available with that scrip")

	}
	var totalQty int = 0
	var totalPrice float64 = 0

	for _, v := range users {
    if strings.ToLower(v.BuysellType) == "buy" {
        totalQty += v.Quantity
        totalPrice += float64(v.Quantity) * v.Price
    } else {
        totalQty -= v.Quantity
        totalPrice -= float64(v.Quantity) * v.Price
    }
}

	fmt.Printf("%s : %d shares, Net Investment: ₹%.2f \n",scrip,totalQty,totalPrice)
	return c.JSON(totalPrice)
}
