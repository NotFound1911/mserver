package recovery

import (
	"github.com/NotFound1911/mserver"
	"net/http"
)

func Recovery() mserver.Middleware {
	return func(ctx *mserver.Context) error {
		defer func() {
			if err := recover(); err != nil {
				ctx.SetStatus(http.StatusInternalServerError).Json(err)
			}
		}()
		ctx.Next()
		return nil
	}
}
