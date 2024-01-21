package recovery

import (
	"github.com/NotFound1911/mserver"
)

type MiddlewareBuilder struct {
	StatusCode int
	Data       []byte
	Log        func(ctx *mserver.Context, err any)
}

func (m MiddlewareBuilder) Build() mserver.Middleware {
	return func(next mserver.HandleFunc) mserver.HandleFunc {
		return func(ctx *mserver.Context) error {
			defer func() {
				if err := recover(); err != nil {
					ctx.SetRespData(m.Data)
					ctx.SetStatus(m.StatusCode)
					m.Log(ctx, err)
				}
			}()
			next(ctx)
			return nil
		}
	}
}
