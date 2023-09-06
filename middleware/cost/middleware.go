package cost

import (
	"github.com/NotFound1911/mserver"
	"log"
	"time"
)

func Cost() mserver.Middleware {
	return func(ctx *mserver.Context) error {
		// 记录开始时间
		start := time.Now()
		log.Printf("api uri start: %v", ctx.GetRequest().RequestURI)
		// 使用next执行具体的业务逻辑
		ctx.Next()
		// 记录结束时间
		end := time.Now()
		cost := end.Sub(start)
		log.Printf("api uri end: %v, cost: %v", ctx.GetRequest().RequestURI, cost.Seconds())
		return nil
	}
}
