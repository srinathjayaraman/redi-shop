package server

import (
	"fmt"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"github.com/valyala/fasthttp"
)

var services = map[string]func(*gorm.DB) fasthttp.RequestHandler{
	"user":    getUserRouter,
	"stock":   getStockRouter,
	"payment": getPaymentRouter,
	"order":   getOrderRouter,
}

// Start initializes the database connection and starts listening to incoming requests
func Start() {
	service := viper.GetString("service")

	// Open database connection
	db, err := gorm.Open("postgres",
		fmt.Sprintf("host=%s port=%d dbname=%s user=%s password=%s sslmode=disable",
			viper.GetString("postgres.url"),
			viper.GetInt("postgres.port"),
			viper.GetString("postgres.database"),
			viper.GetString("postgres.username"),
			viper.GetString("postgres.password"),
		))
	if err != nil {
		logrus.WithError(err).Fatal("unable to connect to database")
	}
	defer db.Close()

	// Get the handlerFunc for the service we want to use
	handlerFn, ok := services[service]
	if !ok {
		logrus.WithField("service", service).Fatal("service does not exist, valid services are: user, stock, order, payment")
	}

	// Start listening to incoming requests
	logrus.WithField("service", service).Info("Redi-shop started, awaiting requests...")
	err = fasthttp.ListenAndServe(":8000", handlerFn(db))
	if err != nil {
		logrus.WithError(err).Fatal("error while listening")
	}
}
