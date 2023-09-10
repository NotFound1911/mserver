package recovery

import (
	"github.com/NotFound1911/mserver"
	"net/http"
)

func Recovery() mserver.Middleware {
	return func(next mserver.HandleFunc) mserver.HandleFunc {
		return func(ctx *mserver.Context) error {
			defer func() {
				if err := recover(); err != nil {
					ctx.SetStatus(http.StatusInternalServerError).Json(err)
				}
			}()
			next(ctx)
			return nil
		}
	}

}
