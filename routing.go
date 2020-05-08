package main

import (
	"fmt"

	"github.com/fasthttp/router"
	"github.com/martijnjanssen/redi-shop/user"
	"github.com/valyala/fasthttp"
)

// returns the router with all user routes
func getUserRouter() fasthttp.RequestHandler {
	h := user.NewRouteHandler()

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

	return r.Handler
}

func getStockRouter() fasthttp.RequestHandler {
	r := router.New()
	r.PanicHandler = panicHandler

	// TODO: Implement

	return r.Handler
}

func getPaymentRouter() fasthttp.RequestHandler {
	r := router.New()
	r.PanicHandler = panicHandler

	// TODO: Implement

	return r.Handler
}

func panicHandler(ctx *fasthttp.RequestCtx, _ interface{}) {
	if r := recover(); r != nil {
		fmt.Println("Recovered in panicHandler", r)
	}

	ctx.Response.Reset()
	ctx.SetStatusCode(fasthttp.StatusInternalServerError)
}
