package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"math"
	"math/rand/v2"
	"trainingmod/database"
	"trainingmod/models"

	"github.com/redis/go-redis/v9"
)

type PaymentHandler struct {
	database.IPaymentDB // prmoted field
		rdb              *redis.Client
	ctx context.Context


}

type IPaymentHandler interface {
	CreatePayment(order models.OrderTable) error
}

func NewPaymentHandler(iPaymentdb database.IPaymentDB, rdb *redis.Client,	ctx context.Context) IPaymentHandler {
	return &PaymentHandler{iPaymentdb,rdb,ctx}
}

func (uh *PaymentHandler) CreatePayment(order models.OrderTable) error {
	Payment := new(models.PaymentTable)
	Payment.OrderId = order.Id

	var err error
	Payment.Amt = 2000 + rand.Float64()*(10000-2000)
	if int(math.Round(Payment.Amt))%2 == 0 {
		Payment.Status = "paid"
	} else {
		Payment.Status = "failed"
	}
	Payment, err = uh.Create(Payment)
	if err != nil {
		return err
	}
	fmt.Println("payment ", Payment)
	bytes,err:=json.Marshal(Payment)
	if err !=nil{
		return err
	}
   err=uh.rdb.Publish(uh.ctx,"payment-success",bytes).Err()
   if err !=nil{
		return err
	}
	println("coming he he")
	return nil

}
