package order

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/go-redis/redis/v8"
	"github.com/gofrs/uuid"
	"github.com/martijnjanssen/redi-shop/util"
	"github.com/sirupsen/logrus"
	"github.com/valyala/fasthttp"
)

type orderStore interface {
	Create(*fasthttp.RequestCtx, string)
	Remove(*fasthttp.RequestCtx, string)
	Find(*fasthttp.RequestCtx, string)
	AddItem(*fasthttp.RequestCtx, string, string)
	RemoveItem(*fasthttp.RequestCtx, string, string)

	GetOrder(*fasthttp.RequestCtx, string) (string, error)
}

var ErrNil = errors.New("value does not exist")

type orderRouteHandler struct {
	orderStore orderStore
	redis      *redis.Client
	chans      map[string]chan (string)
}

func NewRouteHandler(conn *util.Connection) *orderRouteHandler {
	var store orderStore

	switch conn.Backend {
	case util.POSTGRES:
		store = newPostgresOrderStore(conn.Postgres, &conn.URL)
	case util.REDIS:
		store = newRedisOrderStore(conn.Redis, &conn.URL)
	}

	h := &orderRouteHandler{
		orderStore: store,
		redis:      conn.Redis,
		chans:      make(map[string]chan string),
	}

	go handleEvents(conn.Redis, h.chans, util.CHANNEL_ORDER)

	return h
}

func handleEvents(rClient *redis.Client, chans map[string]chan (string), channels ...string) {
	ctx := context.Background()

	pubsub := rClient.Subscribe(ctx, channels...)

	// Wait for confirmation that subscription is created before publishing anything.
	_, err := pubsub.Receive(ctx)
	if err != nil {
		logrus.WithError(err).Panic("error listening to channel")
	}

	var m string
	var rm *redis.Message

	// Go channel which receives messages.
	ch := pubsub.Channel()
	for {
		select {
		case rm = <-ch:
			m = rm.Payload
			s := strings.Split(m, "#")
			chans[s[0]] <- s[1]
		}
	}
}

// Creates order for given user, and returns an order ID
func (h *orderRouteHandler) CreateOrder(ctx *fasthttp.RequestCtx) {
	userID := ctx.UserValue("user_id").(string)
	h.orderStore.Create(ctx, userID)
}

// Deletes an order by ID
func (h *orderRouteHandler) RemoveOrder(ctx *fasthttp.RequestCtx) {
	orderID := ctx.UserValue("order_id").(string)
	h.orderStore.Remove(ctx, orderID)
}

// Retrieves information of an order
func (h *orderRouteHandler) FindOrder(ctx *fasthttp.RequestCtx) {
	orderID := ctx.UserValue("order_id").(string)
	h.orderStore.Find(ctx, orderID)
}

// Adds a g given item in the order given
func (h *orderRouteHandler) AddOrderItem(ctx *fasthttp.RequestCtx) {
	orderID := ctx.UserValue("order_id").(string)
	itemID := ctx.UserValue("item_id").(string)
	h.orderStore.AddItem(ctx, orderID, itemID)
}

// Removes the given item from the give order
func (h *orderRouteHandler) RemoveOrderItem(ctx *fasthttp.RequestCtx) {
	orderID := ctx.UserValue("order_id").(string)
	itemID := ctx.UserValue("item_id").(string)
	h.orderStore.RemoveItem(ctx, orderID, itemID)
}

// Make the payment, subtract the stock and return a status
func (h *orderRouteHandler) CheckoutOrder(ctx *fasthttp.RequestCtx) {
	orderID := ctx.UserValue("order_id").(string)

	// h.orderStore.Checkout(ctx, orderID)
	order, err := h.orderStore.GetOrder(ctx, orderID)
	if err == ErrNil {
		util.NotFound(ctx)
		return
	} else if err != nil {
		logrus.WithError(err).Error("unable to get order")
		util.InternalServerError(ctx)
		return
	}

	trackID := uuid.Must(uuid.NewV4()).String()
	c := make(chan string, 1)

	h.chans[trackID] = c

	// Send message to issue order payment
	err = h.redis.Publish(ctx, util.CHANNEL_PAYMENT, fmt.Sprintf("%s#%s#%s", trackID, util.MESSAGE_PAY, order)).Err()
	if err != nil {
		logrus.WithError(err).Error("unable to publish message")
		util.InternalServerError(ctx)
		return
	}

	select {
	case m := <-c:
		switch m {
		case util.MESSAGE_ORDER_SUCCESS:
			util.Ok(ctx)
		case util.MESSAGE_ORDER_BADREQUEST:
			util.BadRequest(ctx)
		case util.MESSAGE_ORDER_INTERNAL:
			util.InternalServerError(ctx)
		}
	}

	// Close the channel and remove from the map
	close(c)
	delete(h.chans, trackID)
}
