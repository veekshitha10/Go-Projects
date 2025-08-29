package main

import (
	_ "database/sql"
	"flag"
	"fmt"
	"os"
	"time"
	"trainingmod/database"
	"trainingmod/handlers"
	"trainingmod/models"

	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"gorm.io/gorm"
)

var (
	DSN   string
	PORT  string
	debug bool
)

func init() {
	models.Job = make(chan uint, 100)
}
func startJob(db *gorm.DB) {
	go func(){
	for v := range models.Job {
		var order models.OrderTable

		if err := db.First(&order, v).Error; err != nil {
			fmt.Println("Order not found")
			continue
		}
		time.Sleep(time.Second * 2)
		if v%2 == 0 {
			order.Status = "Confirmed"
		} else {
			order.Status = "Failed"

		}
		db.Model(&models.OrderTable{}).Where("id=?", v).Updates(order)
		log.Info().Uint("orderID", v).Str("status", order.Status).Msg("Order processed")

	}
}()
}
func main() {
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
     startJob(db)
	app := fiber.New()
	app.Get("/", handlers.Root)
	app.Get("ping", handlers.Ping)
	app.Get("/health", handlers.Health)

	userHandler := handlers.NewUserHandler(database.NewUserDB(db))
	user_group := app.Group("/api/v1/users")

	user_group.Post("/", userHandler.CreateUser)

	user_group.Get("/:id", userHandler.GetUserBy)

	order_group := app.Group("/api/v1/users/orders")
	order_group.Post("/", userHandler.CreateOrder)
	order_group.Get("/:id", userHandler.GetaOrderBy)
	order_group.Get("/:id/confirm", userHandler.ConfirmOrder)

	app.Listen(":" + PORT)

}

func Init(db *gorm.DB) {
	db.AutoMigrate(&models.UserTable{}, &models.OrderTable{})
}
