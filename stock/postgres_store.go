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
	stock := &Stock{
		Price: price,
	}
	err := s.db.Model(&Stock{}).
		Create(stock).
		Error
	if err != nil {
		util.InternalServerError(ctx)
		return
	}
	response := fmt.Sprintf("{\"id\" : \"%s\"}", stock.ID)
	util.JsonResponse(ctx, fasthttp.StatusCreated, response)
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
	response := fmt.Sprintf("{\"price\" : %d, \"number\": %d}", stock.Price, stock.Number)
	util.JsonResponse(ctx, fasthttp.StatusOK, response)
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
		util.InternalServerError(ctx)
		return
	}

	if stock.Number-number < 0 {
		util.StringResponse(ctx, fasthttp.StatusBadRequest, "failure")
		return
	}

	err = s.db.Model(&Stock{}).
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
