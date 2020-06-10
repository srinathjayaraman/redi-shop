package stock

import (
	"context"
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
	err := s.subtract(ctx, itemID, number)
	if err == util.INTERNAL_ERR {
		util.InternalServerError(ctx)
	} else if err == util.BAD_REQUEST {
		util.BadRequest(ctx)
	}

	util.Ok(ctx)
}

func (s *postgresStockStore) AddStock(ctx *fasthttp.RequestCtx, itemID string, number int) {
	err := s.add(ctx, itemID, number)
	if err == util.INTERNAL_ERR {
		util.InternalServerError(ctx)
	} else if err == util.BAD_REQUEST {
		util.BadRequest(ctx)
	}

	util.Ok(ctx)
}

func (s *postgresStockStore) subtract(_ context.Context, itemID string, number int) error {
	stock := &Stock{}
	err := s.db.Model(&Stock{}).
		Where("id = ?", itemID).
		First(stock).
		Error
	if err == gorm.ErrRecordNotFound {
		return util.BAD_REQUEST
	} else if err != nil {
		logrus.WithError(err).Error("unable to get stock item to subtract")
		return util.INTERNAL_ERR
	}

	if stock.Number-number < 0 {
		logrus.WithField("item_id", itemID).Warning("stock cannot go below 0")
		return util.BAD_REQUEST
	}

	err = s.db.Model(&Stock{}).
		Where("id = ?", itemID).
		Update("number", gorm.Expr("number - ?", number)).
		Error
	if err != nil {
		logrus.WithError(err).Error("unable to subtract stock")
		return util.INTERNAL_ERR
	}

	return nil
}

func (s *postgresStockStore) add(_ context.Context, itemID string, number int) error {
	err := s.db.Model(&Stock{}).
		Where("id = ?", itemID).
		Update("number", gorm.Expr("number + ?", number)).
		Error
	if err == gorm.ErrRecordNotFound {
		return util.BAD_REQUEST
	} else if err != nil {
		logrus.WithError(err).Error("unable to add stock")
		return util.INTERNAL_ERR
	}

	return nil
}
