package user

import (
	"fmt"
	"strconv"

	"github.com/jinzhu/gorm"
	"github.com/valyala/fasthttp"
)

type postgresUserStore struct {
	db *gorm.DB
}

func newPostgresUserStore(db *gorm.DB) *postgresUserStore {
	err := db.AutoMigrate(&User{}).Error
	if err != nil {
		panic(err)
	}

	return &postgresUserStore{
		db: db,
	}
}

func NotFound(ctx *fasthttp.RequestCtx) {
	ctx.SetStatusCode(fasthttp.StatusNotFound)
	ctx.Response.ConnectionClose()
}

func InternalServerError(ctx *fasthttp.RequestCtx) {
	ctx.SetStatusCode(fasthttp.StatusInternalServerError)
	ctx.Response.ConnectionClose()
}

func StringResponse(ctx *fasthttp.RequestCtx, status int, response string) {
	ctx.SetStatusCode(status)
	ctx.SetBodyString(response)
	ctx.Response.ConnectionClose()
}

func (s *postgresUserStore) Create(ctx *fasthttp.RequestCtx) {
	user := &User{}
	err := s.db.Model(&User{}).
		Create(user).
		Error
	if err != nil {
		InternalServerError(ctx)
		return
	}

	StringResponse(ctx, fasthttp.StatusCreated, user.ID)
}

func (s *postgresUserStore) Remove(ctx *fasthttp.RequestCtx, userID string) {
	err := s.db.Model(&User{}).
		Delete(&User{ID: userID}).
		Error
	if err != nil {
		InternalServerError(ctx)
	}

	StringResponse(ctx, fasthttp.StatusOK, "success")
}

func (s *postgresUserStore) Find(ctx *fasthttp.RequestCtx, userID string) {
	user := &User{}
	err := s.db.Model(&User{}).
		Where("id = ?", userID).
		First(user).
		Error
	if err == gorm.ErrRecordNotFound {
		NotFound(ctx)
		return
	} else if err != nil {
		InternalServerError(ctx)
		return
	}

	StringResponse(ctx, fasthttp.StatusOK, fmt.Sprintf("(%s, %d)", user.ID, user.Credit))
}

func (s *postgresUserStore) GetCredit(ctx *fasthttp.RequestCtx, userID string) {
	user := &User{}
	err := s.db.Model(&User{}).
		Where("id = ?", userID).
		First(user).
		Error
	if err == gorm.ErrRecordNotFound {
		NotFound(ctx)
		return
	} else if err != nil {
		InternalServerError(ctx)
		return
	}

	StringResponse(ctx, fasthttp.StatusOK, strconv.Itoa(user.Credit))
}

func (s *postgresUserStore) SubtractCredit(ctx *fasthttp.RequestCtx, userID string, amount int) {
	err := s.db.Model(&User{}).
		Where("id = ?", userID).
		UpdateColumn("credit",
			s.db.Table("users").
				Select("credit - ? as new_credit", amount).
				Where("id = ?", userID).
				SubQuery()).
		Error
	if err != nil {
		StringResponse(ctx, fasthttp.StatusInternalServerError, "failure")
		return
	}

	StringResponse(ctx, fasthttp.StatusOK, "success")
}

func (s *postgresUserStore) AddCredit(ctx *fasthttp.RequestCtx, userID string, amount int) {
	err := s.db.Model(&User{}).
		Where("id = ?", userID).
		UpdateColumn("credit",
			s.db.Table("users").
				Select("credit + ? as new_credit", amount).
				Where("id = ?", userID).
				SubQuery()).
		Error
	if err != nil {
		StringResponse(ctx, fasthttp.StatusInternalServerError, "failure")
		return
	}

	StringResponse(ctx, fasthttp.StatusOK, "success")
}
