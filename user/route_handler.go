package user

import (
	"fmt"
	"strconv"

	"github.com/gofrs/uuid"
	"github.com/valyala/fasthttp"
)

type userRouteHandler struct {
	userStore *userStore
}

func NewRouteHandler() *userRouteHandler {
	return &userRouteHandler{
		userStore: newUserStore(),
	}
}

// Returns an ID for the created user
func (h *userRouteHandler) CreateUser(ctx *fasthttp.RequestCtx) {
	h.userStore.Create(ctx)
}

// Returns success/failure
func (h *userRouteHandler) RemoveUser(ctx *fasthttp.RequestCtx) {
	userID := uuid.Must(uuid.FromString(fmt.Sprintf("%s", ctx.UserValue("user_id"))))

	h.userStore.Remove(ctx, userID)
}

// Returns a user with their details (id, credit)
func (h *userRouteHandler) FindUser(ctx *fasthttp.RequestCtx) {
	userID := uuid.Must(uuid.FromString(fmt.Sprintf("%s", ctx.UserValue("user_id"))))

	h.userStore.Find(ctx, userID)
}

// Returns the current credit of a user
func (h *userRouteHandler) GetUserCredit(ctx *fasthttp.RequestCtx) {
	userID := uuid.Must(uuid.FromString(fmt.Sprintf("%s", ctx.UserValue("user_id"))))

	h.userStore.GetCredit(ctx, userID)
}

// Returns success/failure, depending on the credit status.
// Subtracts the amount from the credit of the user (e.g., to buy an order).
func (h *userRouteHandler) SubtractUserCredit(ctx *fasthttp.RequestCtx) {
	userID := uuid.Must(uuid.FromString(fmt.Sprintf("%s", ctx.UserValue("user_id"))))
	amount, err := strconv.Atoi(fmt.Sprintf("%s", ctx.UserValue("amount")))
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
	userID := uuid.Must(uuid.FromString(fmt.Sprintf("%s", ctx.UserValue("user_id"))))
	amount, err := strconv.Atoi(fmt.Sprintf("%s", ctx.UserValue("amount")))
	if err != nil {
		ctx.SetStatusCode(fasthttp.StatusBadRequest)
		ctx.SetBodyString("amount should be an integer")
		return
	}

	h.userStore.AddCredit(ctx, userID, amount)
}
