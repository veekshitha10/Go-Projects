package handlers

import (
	"context"
	"encoding/json"
	"strconv"

	"mutualfundminiproject/database"
	"mutualfundminiproject/kafka"
	"mutualfundminiproject/models"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/redis/go-redis/v9"
	"github.com/rs/zerolog/log"
)

type UserHandler struct {
	database.IUserDB // prmoted field
	rdb              *redis.Client
	ctx              context.Context
}

type IUserHandler interface {
	CreateUser(c *fiber.Ctx) error
	// GetUserBy(c *fiber.Ctx) error
	CreateOrder(msg *kafka.Messaging) func(c *fiber.Ctx) error // GetaOrderBy(c *fiber.Ctx) error
	//ConfirmOrder(c *fiber.Ctx) error
	GetOrdersByUser(c *fiber.Ctx) error
}

func NewUserHandler(iuserdb database.IUserDB, ctx context.Context, rdb *redis.Client) IUserHandler {
	return &UserHandler{iuserdb, rdb, ctx}
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
	//var str string = ""

	user, err = uh.Create(user)
	if err != nil {
		return err
	}
	log.Logger.Println(user)

	return c.JSON(user)

}

// func (uh *UserHandler) GetUserBy(c *fiber.Ctx) error {
// 	id := c.Params("id") // Retrieves the value of ":id"

// 	_id, err := strconv.Atoi(id)
// 	if err != nil {
// 		return errors.New("invalid id")
// 	}

// 	user, err := uh.GetBy(uint(_id))
// 	if err != nil {
// 		log.Err(err).Msg("data might not be available or some sql issue")
// 		return errors.New("something went wrong or no data available with that id")
// 	}
// 	log.Logger.Println(user)

// 	return c.JSON(user)
// }

// func (uh *UserHandler) GetaOrderBy(c *fiber.Ctx) error {
// 	id := c.Params("id") // Retrieves the value of ":id"

// 	_id, err := strconv.Atoi(id)
// 	if err != nil {
// 		return errors.New("invalid id")
// 	}

// 	user, err := uh.GetOrderBy(uint(_id))
// 	if err != nil {
// 		log.Err(err).Msg("data might not be available or some sql issue")
// 		return errors.New("something went wrong or no data available with that id")
// 	}
// 	log.Logger.Println(user)

// 	return c.JSON(user)
// }

func (uh *UserHandler) CreateOrder(msg *kafka.Messaging) func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		order := new(models.Order)

		err := c.BodyParser(order)
		if err != nil {
			return err
		}
		if _, err := uh.GetOrder(order.SchemeCode); err != nil {
			println("enter proper schemecode")
			return err

		}

		err = order.Validate()
		if err != nil {
			return err
		}
		order.Status = "pending"
		order.PlacedAt = time.Now().Unix()
		nav, naver := uh.rdb.Get(uh.ctx, order.SchemeCode).Float64()
		if naver != nil {
			return naver
		}
		order.NavUsed = nav
		//order, err = uh.IUserDB.CreateOrder(order)
		// if err != nil {
		// 	log.Logger.Println("invalid order request")
		// 	return fiber.NewError(fiber.StatusBadRequest, "invalid order request")
		// }
		log.Logger.Println(order)

		orderBytes, err := order.OrderToBytes()
		if err != nil {
			return err
		}

		msg.ChMessaging <- orderBytes
		v := <-msg.ChMessagingC
		// for v := range msg.ChMessagingC {
		var or models.Order
		er := json.Unmarshal([]byte(v), &or)
		if er != nil {
			return er
		}
		o, err := uh.IUserDB.CreateOrder(&or)
		if err != nil {
			return err
		}
		models.Channel <- o.ID
		print("coming here")

		return c.JSON(o)

		//}
		// er := uh.rdb.Publish(uh.ctx, "user.created", orderBytes).Err()
		// if er != nil {
		// 	return er
		// }

		// paymentREDIS := uh.rdb.Subscribe(uh.ctx, "user.created")
		// Wait for subscription to be created
		// if _, er := paymentREDIS.Receive(uh.ctx); er != nil {
		// 	log.Fatal().Err(er)
		// }

		// ch := paymentREDIS.Channel()
		// for v := range ch {
		// 	var payment models.PaymentTable
		// 	er := json.Unmarshal([]byte(v.Payload), &payment)
		//  if er!=nil{
		// 	log.Fatal().Err(er)

		//  }

		// }
	}
}

// func (h *UserHandler) ConfirmOrder(c *fiber.Ctx) error {
// 	id, err := strconv.Atoi(c.Params("id"))
// 	if err != nil {
// 		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid order id"})
// 	}

// 	_, er := h.GetBy(uint(id))
// 	if er != nil {
// 		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid order id"})
// 	}

//		return c.JSON(fiber.Map{"message": "confirmation started"})
//	}
func (uh *UserHandler) GetOrdersByUser(c *fiber.Ctx) error {
	idParam := c.Params("id")
	userID, err := strconv.Atoi(idParam)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "invalid user id")
	}

	orders, err := uh.IUserDB.GetOrdersByUser(uint(userID))
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	return c.JSON(orders)
}
