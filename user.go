package main

import (
	"fmt"
	"strconv"
	"sync"

	"github.com/gofrs/uuid"
	"github.com/valyala/fasthttp"
)

type store struct {
	// not a concurrent map probably but w/e
	users *sync.Map
	lock  *sync.RWMutex
}

func newStore() *store {
	return &store{
		users: &sync.Map{},
		lock:  &sync.RWMutex{},
	}
}

// Returns an ID for the created user
func (s *store) createUser(ctx *fasthttp.RequestCtx) {
	userID := uuid.Must(uuid.NewV4())

	// If user already exists, return
	if _, ok := s.users.Load(userID); ok {
		ctx.SetStatusCode(fasthttp.StatusInternalServerError)
		return
	}

	s.users.Store(userID, 0)

	ctx.SetBodyString(userID.String())
	ctx.SetStatusCode(fasthttp.StatusCreated)

	// ctx.SetStatusCode(fasthttp.StatusNotImplemented)
}

// Returns success/failure
func (s *store) removeUser(ctx *fasthttp.RequestCtx) {
	userID := uuid.Must(uuid.FromString(fmt.Sprintf("%s", ctx.UserValue("user_id"))))

	// If user does not exist, return
	if _, ok := s.users.Load(userID); !ok {
		ctx.SetStatusCode(fasthttp.StatusNotFound)
		ctx.SetBodyString("failure")
		return
	}

	s.users.Delete(userID)

	ctx.SetStatusCode(fasthttp.StatusOK)
	ctx.SetBodyString("success")

	// ctx.SetStatusCode(fasthttp.StatusNotImplemented)
}

// Returns a user with their details (id, credit)
func (s *store) findUser(ctx *fasthttp.RequestCtx) {
	userID := uuid.Must(uuid.FromString(fmt.Sprintf("%s", ctx.UserValue("user_id"))))

	// Try to get the userCredit
	userCredit, ok := s.users.Load(userID)

	// If user does not exist, return
	if !ok {
		ctx.SetStatusCode(fasthttp.StatusNotFound)
		return
	}

	ctx.SetStatusCode(fasthttp.StatusOK)
	ctx.SetBodyString(fmt.Sprintf("(%s, %d)", userID, userCredit))

	// ctx.SetStatusCode(fasthttp.StatusNotImplemented)
}

// Returns the current credit of a user
func (s *store) getUserCredit(ctx *fasthttp.RequestCtx) {
	userID := uuid.Must(uuid.FromString(fmt.Sprintf("%s", ctx.UserValue("user_id"))))

	// Try to get the user
	userCredit, ok := s.users.Load(userID)

	// If user does not exist, return
	if !ok {
		ctx.SetStatusCode(fasthttp.StatusNotFound)
		return
	}

	ctx.SetStatusCode(fasthttp.StatusOK)
	ctx.SetBodyString(strconv.Itoa(userCredit.(int)))

	// ctx.SetStatusCode(fasthttp.StatusNotImplemented)
}

// Returns success/failure, depending on the credit status.
// Subtracts the amount from the credit of the user (e.g., to buy an order).
func (s *store) subtractUserCredit(ctx *fasthttp.RequestCtx) {
	userID := uuid.Must(uuid.FromString(fmt.Sprintf("%s", ctx.UserValue("user_id"))))
	amount, err := strconv.Atoi(fmt.Sprintf("%s", ctx.UserValue("amount")))
	if err != nil {
		ctx.SetStatusCode(fasthttp.StatusBadRequest)
		ctx.SetBodyString("amount should be an integer")
		return
	}

	// Try to get the user
	userCredit, ok := s.users.Load(userID)

	// If user does not exist, return
	if !ok {
		ctx.SetStatusCode(fasthttp.StatusNotFound)
		ctx.SetBodyString("failure")
		return
	}

	s.users.Store(userID, userCredit.(int)-int(amount))

	ctx.SetStatusCode(fasthttp.StatusOK)
	ctx.SetBodyString("success")

	// ctx.SetStatusCode(fasthttp.StatusNotImplemented)
}

// Returns success/failure, depending on the credit status.
// Adds the amount to the credit of the user.
func (s *store) addUserCredit(ctx *fasthttp.RequestCtx) {
	userID := uuid.Must(uuid.FromString(fmt.Sprintf("%s", ctx.UserValue("user_id"))))
	amount, err := strconv.Atoi(fmt.Sprintf("%s", ctx.UserValue("amount")))
	if err != nil {
		ctx.SetStatusCode(fasthttp.StatusBadRequest)
		ctx.SetBodyString("amount should be an integer")
		return
	}

	// Try to get the user
	userCredit, ok := s.users.Load(userID)

	// If user does not exist, return
	if !ok {
		ctx.SetStatusCode(fasthttp.StatusNotFound)
		ctx.SetBodyString("failure")
		return
	}

	s.users.Store(userID, userCredit.(int)+int(amount))

	ctx.SetStatusCode(fasthttp.StatusOK)
	ctx.SetBodyString("success")

	// ctx.SetStatusCode(fasthttp.StatusNotImplemented)
}
