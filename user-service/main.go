package main

import (
	"os"
	"trainingmod/database"
	"trainingmod/models"
	"trainingmod/handlers"

	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog/log"
	"gorm.io/gorm"
)

var (
	DSN  string
	PORT string
)

func main() {
	service := "users-service"
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
	app := fiber.New()
   
	app.Get("/", handlers.Root)
	app.Get("ping", handlers.Ping)
	app.Get("/health", handlers.Health)

	userHandler := handlers.NewOrderHandler(database.NewOrderDB(db))
	user_group := app.Group("/api/v1/placeorders")

	user_group.Post("/", userHandler.CreateOrder)
				user_group.Get("/:scrip", userHandler.CalculateNet)

	 log.Fatal().
        Err(app.Listen(":" + PORT)).
        Msg("Fiber server failed to start")

}

func Init(db *gorm.DB) {
	db.AutoMigrate(&models.OrdersModel{})
}
