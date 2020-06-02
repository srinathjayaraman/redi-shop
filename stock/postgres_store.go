package stock

import (
	"fmt"

	"github.com/jinzhu/gorm"
	"github.com/martijnjanssen/redi-shop/util"
	"github.com/sirupsen/logrus"
	"github.com/valyala/fasthttp"
)

type postgresStockStore struct {
	db   *gorm.DB
	urls *util.Services
}

func newPostgresStockStore(db *gorm.DB, urls *util.Services) *postgresStockStore {
	// AutoMigrate structs to create or update database tables
	err := db.AutoMigrate(&Stock{}).Error
	if err != nil {
		panic(err)
	}

	return &postgresStockStore{
		db:   db,
		urls: urls,
	}
}

func (s *postgresStockStore) Create(ctx *fasthttp.RequestCtx, price int) {
	stock := &Stock{
		Price: price,
	}
	err := s.db.Model(&Stock{}).
		Create(stock).
		Error
	if err != nil {
		logrus.WithError(err).Error("unable to create new stock item")
		util.InternalServerError(ctx)
		return
	}

	response := fmt.Sprintf("{\"item_id\": \"%s\"}", stock.ID)
	util.JSONResponse(ctx, fasthttp.StatusCreated, response)
}

func (s *postgresStockStore) Find(ctx *fasthttp.RequestCtx, itemID string) {
	stock := &Stock{}
	err := s.db.Model(&Stock{}).
		Where("id = ?", itemID).
		First(stock).
		Error
	if err == gorm.ErrRecordNotFound {
		util.NotFound(ctx)
		return
	} else if err != nil {
		logrus.WithError(err).Error("unable to find stock item")
		util.InternalServerError(ctx)
		return
	}

	response := fmt.Sprintf("{\"stock\": %d, \"price\": %d}", stock.Number, stock.Price)
	util.JSONResponse(ctx, fasthttp.StatusOK, response)
}

func (s *postgresStockStore) SubtractStock(ctx *fasthttp.RequestCtx, itemID string, number int) {
	stock := &Stock{}
	err := s.db.Model(&Stock{}).
		Where("id = ?", itemID).
		First(stock).
		Error
	if err == gorm.ErrRecordNotFound {
		util.NotFound(ctx)
		return
	} else if err != nil {
		logrus.WithError(err).Error("unable to get stock item to subtract")
		util.InternalServerError(ctx)
		return
	}

	if stock.Number-number < 0 {
		logrus.WithField("item_id", itemID).Warning("stock cannot go below 0")
		util.BadRequest(ctx)
		return
	}

	err = s.db.Model(&Stock{}).
		Where("id = ?", itemID).
		Update("number", gorm.Expr("number - ?", number)).
		Error
	if err != nil {
		logrus.WithError(err).Error("unable to subtract stock")
		util.InternalServerError(ctx)
		return
	}

	util.Ok(ctx)
}

func (s *postgresStockStore) AddStock(ctx *fasthttp.RequestCtx, itemID string, number int) {
	err := s.db.Model(&Stock{}).
		Where("id = ?", itemID).
		Update("number", gorm.Expr("number + ?", number)).
		Error
	if err == gorm.ErrRecordNotFound {
		util.NotFound(ctx)
		return
	} else if err != nil {
		util.InternalServerError(ctx)
		logrus.WithError(err).Error("unable to add stock")
		return
	}

	util.Ok(ctx)
}
