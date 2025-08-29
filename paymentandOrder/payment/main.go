package main

import (
	"context"
	_ "database/sql"
	"encoding/json"
	"flag"
	"os"
	"time"
	"trainingmod/database"
	"trainingmod/handlers"
	"trainingmod/models"

	"github.com/redis/go-redis/v9"
	_ "github.com/redis/go-redis/v9"

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

func main() {
	service := "Payments-service"
	flag.BoolVar(&debug, "debug", false, "sets log level to debug")
	flag.Parse()
	zerolog.SetGlobalLevel(zerolog.InfoLevel)
	if debug {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	}

	DSN = os.Getenv("DSN")
	if DSN == "" {
		DSN = `host=localhost user=app password=app123 dbname=paymentdb port=5432 sslmode=disable`
		log.Info().Msg(DSN)
	}
	PORT = os.Getenv("PORT")
	if PORT == "" {
		PORT = "8089"
	}

	db, err := database.GetConnection(DSN)
	ctx := context.Background()

	if err != nil {
		//log.Fatal().Msg("unable to connect to the database..." + err.Error())
		log.Fatal().
			Err(err).
			Str("service", service).
			Msgf("unable to connect to the database %s", service)
	}
	log.Info().Str("service", service).Msg("database connection is established")
	Init(db)
	rdb := redis.NewClient(
		&redis.Options{
			Addr: "localhost:6379", 
			
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
	PaymentHandler := handlers.NewPaymentHandler(database.NewPaymentDB(db),rdb,ctx)
	
	paymentREDIS := rdb.Subscribe(ctx, "user.created")
	// Wait for subscription to be created
	if _, er := paymentREDIS.Receive(ctx); er != nil {
		log.Fatal().Err(er)
	}

	ch := paymentREDIS.Channel()
	for v := range ch {
		var order models.OrderTable
		er := json.Unmarshal([]byte(v.Payload), &order)
	 if er!=nil{
		log.Fatal().Err(er)

	 }
	 PaymentHandler.CreatePayment(order)
	}

}

func Init(db *gorm.DB) {
	db.AutoMigrate(&models.PaymentTable{})
}
