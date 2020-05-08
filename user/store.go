package user

import (
	"fmt"
	"strconv"
	"sync"

	"github.com/gofrs/uuid"
	"github.com/valyala/fasthttp"
)

type userStore struct {
	users *sync.Map
}

func newUserStore() *userStore {
	return &userStore{
		users: &sync.Map{},
	}
}

func (s *userStore) Create(ctx *fasthttp.RequestCtx) {
	userID := uuid.Must(uuid.NewV4())

	// If user already exists, return
	if _, ok := s.users.Load(userID); ok {
		ctx.SetStatusCode(fasthttp.StatusInternalServerError)
		return
	}

	s.users.Store(userID, 0)

	ctx.SetBodyString(userID.String())
	ctx.SetStatusCode(fasthttp.StatusCreated)
}

func (s *userStore) Remove(ctx *fasthttp.RequestCtx, userID uuid.UUID) {
	// If user does not exist, return
	if _, ok := s.users.Load(userID); !ok {
		ctx.SetStatusCode(fasthttp.StatusNotFound)
		ctx.SetBodyString("failure")
		return
	}

	s.users.Delete(userID)

	ctx.SetStatusCode(fasthttp.StatusOK)
	ctx.SetBodyString("success")
}

func (s *userStore) Find(ctx *fasthttp.RequestCtx, userID uuid.UUID) {
	// Try to get the userCredit
	userCredit, ok := s.users.Load(userID)

	// If user does not exist, return
	if !ok {
		ctx.SetStatusCode(fasthttp.StatusNotFound)
		return
	}

	ctx.SetStatusCode(fasthttp.StatusOK)
	ctx.SetBodyString(fmt.Sprintf("(%s, %d)", userID, userCredit))
}

func (s *userStore) GetCredit(ctx *fasthttp.RequestCtx, userID uuid.UUID) {
	// Try to get the user
	userCredit, ok := s.users.Load(userID)

	// If user does not exist, return
	if !ok {
		ctx.SetStatusCode(fasthttp.StatusNotFound)
		return
	}

	ctx.SetStatusCode(fasthttp.StatusOK)
	ctx.SetBodyString(strconv.Itoa(userCredit.(int)))
}

func (s *userStore) SubtractCredit(ctx *fasthttp.RequestCtx, userID uuid.UUID, amount int) {
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
}

func (s *userStore) AddCredit(ctx *fasthttp.RequestCtx, userID uuid.UUID, amount int) {
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
}
