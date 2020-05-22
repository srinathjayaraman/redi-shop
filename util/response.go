package util

import (
	"github.com/valyala/fasthttp"
)

func NotFound(ctx *fasthttp.RequestCtx) {
	ctx.SetStatusCode(fasthttp.StatusNotFound)
}

func InternalServerError(ctx *fasthttp.RequestCtx) {
	ctx.SetStatusCode(fasthttp.StatusInternalServerError)
}

func StringResponse(ctx *fasthttp.RequestCtx, status int, response string) {
	ctx.SetStatusCode(status)
	ctx.SetBodyString(response)
}

func JsonResponse(ctx *fasthttp.RequestCtx, status int, response string) {
	ctx.SetStatusCode(status)
	ctx.SetBodyString(response)
	ctx.SetContentType("application/json")
}
