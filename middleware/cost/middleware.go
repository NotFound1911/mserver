package cost

import (
	"github.com/NotFound1911/mserver"
	"log"
	"time"
)

type MiddlewareBuilder struct {
}

func (m MiddlewareBuilder) Build() mserver.Middleware {
	return func(next mserver.HandleFunc) mserver.HandleFunc {
		return func(ctx *mserver.Context) error {
			// 记录开始时间
			start := time.Now()
			log.Printf("api uri start: %v", ctx.GetRequest().RequestURI)
			// 使用handleFunc执行具体的业务逻辑
			next(ctx)
			// 记录结束时间
			end := time.Now()
			cost := end.Sub(start)
			log.Printf("api uri end: %v, cost: %v", ctx.GetRequest().RequestURI, cost.Seconds())
			return nil
		}
	}
}
