package util

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/sirupsen/logrus"
	"github.com/valyala/fasthttp"
)

var services = map[string]func(*gorm.DB) fasthttp.RequestHandler{
	"user":    getUserRouter,
	"stock":   getStockRouter,
	"payment": getPaymentRouter,
	"order":   getOrderRouter,
}

func StartServer(serviceName string) {
	db, err := gorm.Open("postgres", "host=localhost port=5432 user=postgres dbname=redi password=postgres sslmode=disable")
	if err != nil {
		logrus.WithError(err).Fatal("unable to connect to database")
	}
	defer db.Close()

	handlerFn, ok := services[serviceName]
	if !ok {
		logrus.WithField("service", serviceName).Fatal("service does not exist, valid services are: user, stock, order, payment")
	}

	logrus.WithField("service", serviceName).Info("Redi-shop started, awaiting requests...")
	err = fasthttp.ListenAndServe(":8000", handlerFn(db))
	if err != nil {
		logrus.WithError(err).Fatal("error while listening")
	}
}
