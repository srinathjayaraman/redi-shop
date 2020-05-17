package util

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/sirupsen/logrus"
	"github.com/valyala/fasthttp"
)

func StartServer(serverName string, getRouterFunc func(*gorm.DB) fasthttp.RequestHandler) {
	db, err := gorm.Open("postgres", "host=localhost port=5432 user=postgres dbname=redi password=postgres sslmode=disable")
	if err != nil {
		logrus.WithError(err).Fatal("unable to connect to database")
	}
	defer db.Close()

	logrus.WithField("service", serverName).Info("Redi-shop started, awaiting requests...")
	err = fasthttp.ListenAndServe(":8000", getRouterFunc(db))
	if err != nil {
		logrus.WithError(err).Fatal("error while listening")
	}
}
