package main

import (
	"fmt"

	"github.com/fasthttp/router"
	"github.com/jinzhu/gorm"
	"github.com/martijnjanssen/redi-shop/user"
	"github.com/valyala/fasthttp"
)

// returns the router with all user routes
func getUserRouter(db *gorm.DB) fasthttp.RequestHandler {
	h := user.NewRouteHandler(db)

	r := router.New()
	r.PanicHandler = panicHandler

	r.POST("/users/create/", h.CreateUser)
	r.DELETE("/users/remove/{user_id}", h.RemoveUser)
	r.GET("/users/find/{user_id}", h.FindUser)

	r.GET("/users/credit/{user_id}", h.GetUserCredit)
	r.POST("/users/credit/subtract/{user_id}/{amount}", h.SubtractUserCredit)
	r.POST("/users/credit/add/{user_id}/{amount}", h.AddUserCredit)

	return r.Handler
}

func getOrderRouter() fasthttp.RequestHandler {
	r := router.New()
	r.PanicHandler = panicHandler

	// TODO: Implement
	r.POST("/orders/create/{user_id}", nil)
	r.DELETE("/orders/remove/{order_id}", nil)
	r.GET("/orders/find/{order_id}", nil)
	r.POST("/orders/additem/{order_id}/{item_id}", nil)
	r.DELETE("/orders/removeitem/{order_id}/{item_id}", nil)
	r.POST("/orders/checkout/{order_id}", nil)

	return r.Handler
}

func getStockRouter() fasthttp.RequestHandler {
	r := router.New()
	r.PanicHandler = panicHandler

	// TODO: Implement
	r.GET("/stock/availability/{item_id}", nil)
	r.POST("/stock/subtract/{item_id}/{number}", nil)
	r.POST("/stock/add/{item_id}/{number}", nil)
	// TODO: This route is probably incorrectly specified in the doc, should probably have a price
	r.POST("/stock/item/create", nil)

	return r.Handler
}

func getPaymentRouter() fasthttp.RequestHandler {
	r := router.New()
	r.PanicHandler = panicHandler

	// TODO: Implement
	r.POST("/payment/pay/{user_id}/{order_id}", nil)
	r.POST("/payment/cancel/{user_id}/{order_id}", nil)
	r.GET("/payment/status/{order_id}", nil)

	return r.Handler
}

func panicHandler(ctx *fasthttp.RequestCtx, _ interface{}) {
	if r := recover(); r != nil {
		fmt.Println("Recovered in panicHandler", r)
	}

	ctx.Response.Reset()
	ctx.SetStatusCode(fasthttp.StatusInternalServerError)
}
