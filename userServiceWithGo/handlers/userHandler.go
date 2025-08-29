package handlers

import (
	"errors"
	"strconv"
	"time"
	"trainingmod/database"
	"trainingmod/models"

	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog/log"
)

type UserHandler struct {
	database.IUserDB // prmoted field
}

type IUserHandler interface {
	CreateUser(c *fiber.Ctx) error
	GetUserBy(c *fiber.Ctx) error
	CreateOrder(c *fiber.Ctx) error
		GetaOrderBy(c *fiber.Ctx) error
		ConfirmOrder(c *fiber.Ctx) error

}

func NewUserHandler(iuserdb database.IUserDB) IUserHandler {
	return &UserHandler{iuserdb}
}

func (uh *UserHandler) CreateUser(c *fiber.Ctx) error {
	user := new(models.UserTable)
	err := c.BodyParser(user)
	if err != nil {
		return err
	}

	err = user.Validate()
	if err != nil {
		return err
	}

	user.LastModified = time.Now().Unix()

	user, err = uh.Create(user)
	if err != nil {
		return err
	}
				log.Logger.Println(user)

	return c.JSON(user)

}

func (uh *UserHandler) GetUserBy(c *fiber.Ctx) error {
	id := c.Params("id") // Retrieves the value of ":id"

	_id, err := strconv.Atoi(id)
	if err != nil {
		return errors.New("invalid id")
	}

	user, err := uh.GetBy(uint(_id))
	if err != nil {
		log.Err(err).Msg("data might not be available or some sql issue")
		return errors.New("something went wrong or no data available with that id")
	}
			log.Logger.Println(user)

	return c.JSON(user)
}

func (uh *UserHandler) GetaOrderBy(c *fiber.Ctx) error {
	id := c.Params("id") // Retrieves the value of ":id"

	_id, err := strconv.Atoi(id)
	if err != nil {
		return errors.New("invalid id")
	}

	user, err := uh.GetOrderBy(uint(_id))
	if err != nil {
		log.Err(err).Msg("data might not be available or some sql issue")
		return errors.New("something went wrong or no data available with that id")
	}
			log.Logger.Println(user)

	return c.JSON(user)
}


func (uh *UserHandler) CreateOrder(c *fiber.Ctx) error {
	order := new(models.OrderTable)
	err := c.BodyParser(order)
	if err != nil {
		return err
	}
	err = order.Validate()
	if err != nil {
		return err
	}
	   order.Status="pending"


	order.LastModified = time.Now().Unix()

	order, err = uh.IUserDB.CreateOrder(order)
	if err != nil {
		log.Logger.Println("invalid order request")
		return fiber.NewError(fiber.StatusBadRequest, "invalid order request")
	}
			log.Logger.Println(order)

	return c.JSON(order)
}
func (h *UserHandler) ConfirmOrder(c *fiber.Ctx) error {
    id, err := strconv.Atoi(c.Params("id"))
    if err != nil {
        return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid order id"})
    }

    	_, er := h.GetBy(uint(id))
if er != nil {
        return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid order id"})
    }

   models.Job<-uint( id)
    return c.JSON(fiber.Map{"message": "confirmation started"})
}
