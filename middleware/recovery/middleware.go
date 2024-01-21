package recovery

import (
	"github.com/NotFound1911/mserver"
	"net/http"
)

type MiddlewareBuilder struct {
}

func (m MiddlewareBuilder) Build() mserver.Middleware {
	return func(handleFunc mserver.HandleFunc) mserver.HandleFunc {
		return func(ctx *mserver.Context) error {
			defer func() {
				if err := recover(); err != nil {
					ctx.SetStatus(http.StatusInternalServerError).Json(err)
				}
			}()
			handleFunc(ctx)
			return nil
		}
	}
}
