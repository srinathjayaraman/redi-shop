package user

import (
	"strconv"

	"github.com/jinzhu/gorm"
	"github.com/valyala/fasthttp"
)

type userStore interface {
	Create(*fasthttp.RequestCtx)
	Remove(*fasthttp.RequestCtx, string)
	Find(*fasthttp.RequestCtx, string)
	AddCredit(*fasthttp.RequestCtx, string, int)
	SubtractCredit(*fasthttp.RequestCtx, string, int)
}

type userRouteHandler struct {
	userStore userStore
}

func NewRouteHandler(db *gorm.DB) *userRouteHandler {
	return &userRouteHandler{
		userStore: newPostgresUserStore(db),
	}
}

// Returns an ID for the created user
func (h *userRouteHandler) CreateUser(ctx *fasthttp.RequestCtx) {
	h.userStore.Create(ctx)
}

// Returns success/failure
func (h *userRouteHandler) RemoveUser(ctx *fasthttp.RequestCtx) {
	userID := ctx.UserValue("user_id").(string)

	h.userStore.Remove(ctx, userID)
}

// Returns a user with their details (id, credit)
func (h *userRouteHandler) FindUser(ctx *fasthttp.RequestCtx) {
	userID := ctx.UserValue("user_id").(string)

	h.userStore.Find(ctx, userID)
}

// Returns success/failure, depending on the credit status.
// Subtracts the amount from the credit of the user (e.g., to buy an order).
func (h *userRouteHandler) SubtractUserCredit(ctx *fasthttp.RequestCtx) {
	userID := ctx.UserValue("user_id").(string)
	amount, err := strconv.Atoi(ctx.UserValue("amount").(string))
	if err != nil {
		ctx.SetStatusCode(fasthttp.StatusBadRequest)
		ctx.SetBodyString("amount should be an integer")
		return
	}

	h.userStore.SubtractCredit(ctx, userID, amount)
}

// Returns success/failure, depending on the credit status.
// Adds the amount to the credit of the user.
func (h *userRouteHandler) AddUserCredit(ctx *fasthttp.RequestCtx) {
	userID := ctx.UserValue("user_id").(string)
	amount, err := strconv.Atoi(ctx.UserValue("amount").(string))
	if err != nil {
		ctx.SetStatusCode(fasthttp.StatusBadRequest)
		ctx.SetBodyString("amount should be an integer")
		return
	}

	h.userStore.AddCredit(ctx, userID, amount)
}
