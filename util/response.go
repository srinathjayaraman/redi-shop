package util

import (
	"github.com/valyala/fasthttp"
)

func Ok(ctx *fasthttp.RequestCtx) {
	ctx.SetStatusCode(fasthttp.StatusOK)
}

func NotFound(ctx *fasthttp.RequestCtx) {
	ctx.SetStatusCode(fasthttp.StatusNotFound)
}

func BadRequest(ctx *fasthttp.RequestCtx) {
	ctx.SetStatusCode(fasthttp.StatusBadRequest)
}

func InternalServerError(ctx *fasthttp.RequestCtx) {
	ctx.SetStatusCode(fasthttp.StatusInternalServerError)
}

func JSONResponse(ctx *fasthttp.RequestCtx, status int, response string) {
	ctx.SetStatusCode(status)
	ctx.SetBodyString(response)
	ctx.SetContentType("application/json")
}
