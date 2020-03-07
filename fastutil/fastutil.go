package fastutil

import (
	"strconv"

	"github.com/valyala/fasthttp"
	"github.com/wongoo/webu/cerror"
)

func RequireQueryArg(ctx *fasthttp.RequestCtx, arg string) ([]byte, error) {
	val := ctx.QueryArgs().Peek(arg)
	if len(val) == 0 {
		return nil, cerror.ErrArgRequired
	}

	return val, nil
}

func RequireQueryInt(ctx *fasthttp.RequestCtx, arg string) (int, error) {
	val, err := RequireQueryArg(ctx, arg)
	if err != nil {
		return 0, err
	}

	i, err := strconv.Atoi(string(val))
	if err != nil {
		return 0, cerror.ErrValueInvalid
	}

	return i, nil
}
