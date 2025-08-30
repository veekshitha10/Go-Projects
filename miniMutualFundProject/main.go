package main

import (
	"context"
	_ "database/sql"
	"flag"
	"fmt"
	"math/rand"
	"mutualfundminiproject/database"
	"mutualfundminiproject/handlers"
	"mutualfundminiproject/kafka"
	"mutualfundminiproject/models"
	"os"
	"strings"
	"time"

	"github.com/ansrivas/fiberprometheus/v2"
	"github.com/gofiber/fiber/v2"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"github.com/redis/go-redis/v9"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"gorm.io/gorm"
)

var (
	DSN     string
	PORT    string
	debug   bool
	schemes []models.SchemeTable
)

func startJob(db *gorm.DB) {
	go func() {

		for v := range models.Channel {
			print("coming hereffvfvfv")

			var order models.Order

			if err := db.First(&order, v).Error; err != nil {
				fmt.Println("Order not found")
				continue
			}
			time.Sleep(time.Second * 15)
			order.Status = "placed"
			order.ConfirmedAt = time.Now().Unix()

			//.Sleep(time.Second * 2)

			db.Model(&models.Order{}).Where("id=?", v).Updates(order)
			var holding models.HoldingsTable

			if err := db.First(&holding, order.SchemeCode).Error; err != nil {
				s := models.HoldingsTable{
					UserId:     order.UserId,
					SchemeCode: order.SchemeCode,
					Units:      order.Units,
				}
				if err := db.Create(&s).Error; err != nil {
					panic(fmt.Sprint(err))
				}
			}
			var units float64 = holding.Units
			if order.Side == "buy" {
				units = holding.Units + order.Units
			} else if order.Side == "sell" || order.Units > 0 {
				units = holding.Units - order.Units

			}
			if units < 0 {
				units = 0
			}
			holding.Units = units
			db.Model(&models.HoldingsTable{}).Where("scheme_code =?", order.SchemeCode).Updates(holding)

		}
	}()
}
func updateRedisValue(rdb *redis.Client, ctx context.Context) {
	if len(schemes) > 0 {
		for _, v := range schemes {

			min := 1.0
			max := 10000.0

			randomFloat := min + rand.Float64()*(max-min)
			err := rdb.Set(ctx, v.SchemeCode, randomFloat, 0).Err()
			if err != nil {
				panic(fmt.Sprint(err))
			}

		}
		time.Sleep(time.Second * 5)
	}
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
	SEEDS := os.Getenv("KAFKA_BROKERS")
	if SEEDS == "" {
		SEEDS = "localhost:19092, localhost:29092, localhost:39092"
	}
	log.Info().Str("service", service).Msg("database connection is established")
	Init(db)
	s3Client, err := minio.New(models.Endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(models.AccessKey, models.SecretKey, ""),
		Secure: false,
		Region: models.Region, // optional but nice to set
	})

	userHandler := handlers.NewUserHandler(database.NewUserDB(db), &ctx, rdb, s3Client)
	go updateRedisValue(rdb, ctx)
	startJob(db)
	msgUsersCreated := kafka.NewMessaging("omnenext.users.created", strings.Split(SEEDS, ","), "--cg=demo-consumer-group")
	go msgUsersCreated.ProduceRecords()
	go msgUsersCreated.ConsumeRecords()

	if err != nil {
		log.Fatal().
			Err(err)
	}

	app := fiber.New(
		fiber.Config{
			BodyLimit: 50 * 1024 * 1024})
	prom := fiberprometheus.New(service)
	prom.RegisterAt(app, "/metrics") // exposes Prometheus metrics here
	app.Use(prom.Middleware)
	app.Get("/", handlers.Root)
	app.Get("ping", handlers.Ping)
	app.Get("/health", handlers.Health)

	user_group := app.Group("/api/v1/users")

	user_group.Post("/login", userHandler.CreateUser)
	user_group.Post("/UploadImg/:id", userHandler.UploadImage)

	user_group.Get("/:id", userHandler.GetOrdersByUser)

	order_group := app.Group("/api/v1/users/orders")
	order_group.Post("/", userHandler.CreateOrder(msgUsersCreated))
	// order_group.Get("/:id", userHandler.GetaOrderBy)
	//	order_group.Get("/:id/confirm", userHandler.ConfirmOrder)

	app.Listen(":" + PORT)

}

func Init(db *gorm.DB) {
	db.AutoMigrate(&models.UserTable{}, &models.Order{}, &models.HoldingsTable{}, &models.SchemeTable{})
	if err := models.SeedSchemeTableIfEmpty(db); err != nil {
		panic(fmt.Sprint(err))
	}
	err := db.Find(&schemes).Error
	if err != nil {
		// handle error
		fmt.Println("Error fetching schemes:", err)
	} else {
		fmt.Println("Schemes:", schemes)
	}
	models.Channel = make(chan uint)
}
