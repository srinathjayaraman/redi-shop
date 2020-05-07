package main

import (
	"github.com/fasthttp/router"
	"github.com/valyala/fasthttp"
)

// return the router with all routes
func getRouter() fasthttp.RequestHandler {
	r := router.New()

	store := newStore()

	r.POST("/users/create/", store.createUser)
	r.DELETE("/users/remove/{user_id}", store.removeUser)
	r.GET("/users/find/{user_id}", store.findUser)

	r.GET("/users/credit/{user_id}", store.getUserCredit)
	r.POST("/users/credit/subtract/{user_id}/{amount}", store.subtractUserCredit)
	r.POST("/users/credit/add/{user_id}/{amount}", store.addUserCredit)

	r.PanicHandler = func(ctx *fasthttp.RequestCtx, body interface{}) {
		ctx.Response.Reset()
		ctx.SetStatusCode(fasthttp.StatusInternalServerError)
	}

	return r.Handler
}
