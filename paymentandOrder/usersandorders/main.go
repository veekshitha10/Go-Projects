package main

import (
	"context"
	_ "database/sql"
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"time"
	"trainingmod/database"
	"trainingmod/handlers"
	"trainingmod/models"

	"github.com/gofiber/fiber/v2"
	"github.com/redis/go-redis/v9"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"gorm.io/gorm"
)

var (
	DSN   string
	PORT  string
	debug bool
)

func startJob(db *gorm.DB) {
	go func() {
		for v := range models.Channel {
			var payment models.PaymentTable
			er := json.Unmarshal([]byte(v.Payload), &payment)
			if er != nil {
				log.Fatal().Err(er)

			}
			var order models.OrderTable

			if err := db.First(&order, payment.OrderId).Error; err != nil {
				fmt.Println("Order not found")
				continue
			}
			order.Orderstatus = payment.Status
			//.Sleep(time.Second * 2)

			db.Model(&models.OrderTable{}).Where("id=?", payment.OrderId).Updates(order)

		}
	}()
}
func main() {
	ctx := context.Background()

	service := "users-service"
	flag.BoolVar(&debug, "debug", false, "sets log level to debug")
	flag.Parse()
	zerolog.SetGlobalLevel(zerolog.InfoLevel)
	if debug {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	}

	DSN = os.Getenv("DSN")
	if DSN == "" {
		DSN = `host=localhost user=app password=app123 dbname=usersdb port=5432 sslmode=disable`
		log.Info().Msg(DSN)
	}
	PORT = os.Getenv("PORT")
	if PORT == "" {
		PORT = "8089"
	}
	rdb := redis.NewClient(
		&redis.Options{
			Addr: "localhost:6379", 
			// Password:  "",               
			// DB:        0,
			DialTimeout:  2 * time.Second,
			ReadTimeout:  1 * time.Second,
			WriteTimeout: 1 * time.Second,
			PoolSize:     20,
			MinIdleConns: 4,
		})
	defer rdb.Close()

	if err := rdb.Ping(ctx).Err(); err != nil {
		log.Fatal().
			Err(err)
	} else {
		println("Connected")
	}
	db, err := database.GetConnection(DSN)
	if err != nil {
		//log.Fatal().Msg("unable to connect to the database..." + err.Error())
		log.Fatal().
			Err(err).
			Str("service", service).
			Msgf("unable to connect to the database %s", service)
	}

	log.Info().Str("service", service).Msg("database connection is established")
	Init(db)
	userHandler := handlers.NewUserHandler(database.NewUserDB(db), rdb, ctx)

	res := rdb.Subscribe(ctx, "payment-success")
	print("cominge here")
	if _, er := res.Receive(ctx); er != nil {
		print("cominge here error")

		log.Fatal().Err(er)
	}
	models.Channel = res.Channel()
	print("cominge 2")
	startJob(db)

	app := fiber.New()
	app.Get("/", handlers.Root)
	app.Get("ping", handlers.Ping)
	app.Get("/health", handlers.Health)

	user_group := app.Group("/api/v1/users")

	user_group.Post("/", userHandler.CreateUser)

	user_group.Get("/:id", userHandler.GetUserBy)

	order_group := app.Group("/api/v1/users/orders")
	order_group.Post("/", userHandler.CreateOrder)
	order_group.Get("/:id", userHandler.GetaOrderBy)
	//	order_group.Get("/:id/confirm", userHandler.ConfirmOrder)

	app.Listen(":" + PORT)

}

func Init(db *gorm.DB) {
	db.AutoMigrate(&models.UserTable{}, &models.OrderTable{})
}
