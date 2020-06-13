package stock

import (
	"context"
	"fmt"

	"github.com/jinzhu/gorm"
	"github.com/martijnjanssen/redi-shop/util"
	"github.com/pkg/errors"
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
	var result error

	err := s.db.Transaction(func(tx *gorm.DB) error {
		err := tx.Model(&Stock{}).
			Where("id = ?", itemID).
			Update("number", gorm.Expr("number - ?", number)).
			Error
		if err == gorm.ErrRecordNotFound {
			result = util.BAD_REQUEST
			return errors.Wrap(err, "stock not found")
		} else if err != nil {
			result = util.INTERNAL_ERR
			return errors.Wrap(err, "unable to subtract stock")
		}

		// Check whether the stock is still above 0
		err = tx.Model(&Stock{}).
			Where("id = ?", itemID).
			Where("number-? > 0", number).
			First(&Stock{}).
			Error
		if err == gorm.ErrRecordNotFound {
			result = util.BAD_REQUEST
			return errors.Wrap(err, "stock cannot go below 0")
		} else if err != nil {
			result = util.INTERNAL_ERR
			return errors.Wrap(err, "unable to get stock")
		}

		return nil
	})
	if err != nil {
		logrus.WithError(err).Error("unable to subtract stock")
		return result
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
