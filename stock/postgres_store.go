package stock

import (
	"fmt"

	"github.com/jinzhu/gorm"
	"github.com/martijnjanssen/redi-shop/util"
	"github.com/valyala/fasthttp"
)

type postgresStockStore struct {
	db *gorm.DB
}

func newPostgresStockStore(db *gorm.DB) *postgresStockStore {
	// AutoMigrate structs to create or update database tables
	err := db.AutoMigrate(&Stock{}).Error
	if err != nil {
		panic(err)
	}

	return &postgresStockStore{
		db: db,
	}
}

func (s *postgresStockStore) Create(ctx *fasthttp.RequestCtx, price int) {
	stock := &Stock{}
	err := s.db.Model(&Stock{}).
		Create(stock).
		Update("price", price).
		Error
	if err != nil {
		util.InternalServerError(ctx)
		return
	}

	util.StringResponse(ctx, fasthttp.StatusCreated, stock.ID)
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
		util.InternalServerError(ctx)
		return
	}

	util.StringResponse(ctx, fasthttp.StatusOK, fmt.Sprintf("(%d, %d)", stock.Price, stock.Number))
}

func (s *postgresStockStore) SubtractStock(ctx *fasthttp.RequestCtx, itemID string, number int) {
	err := s.db.Model(&Stock{}).
		Where("id = ?", itemID).
		Update("number", gorm.Expr("number - ?", number)).
		Error

	if err != nil {
		util.StringResponse(ctx, fasthttp.StatusInternalServerError, "failure")
		return
	}

	util.StringResponse(ctx, fasthttp.StatusOK, "success")
}

func (s *postgresStockStore) AddStock(ctx *fasthttp.RequestCtx, itemID string, number int) {
	err := s.db.Model(&Stock{}).
		Where("id = ?", itemID).
		Update("number", gorm.Expr("number + ?", number)).
		Error
	if err != nil {
		util.StringResponse(ctx, fasthttp.StatusInternalServerError, "failure")
		return
	}

	util.StringResponse(ctx, fasthttp.StatusOK, "success")
}
