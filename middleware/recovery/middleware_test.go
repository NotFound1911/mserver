package recovery

import (
	"fmt"
	"github.com/NotFound1911/mserver"
	"testing"
)

func TestMiddlewareBuilder_Build(t *testing.T) {
	builder := MiddlewareBuilder{
		StatusCode: 500,
		Data:       []byte("发生 panic 了"),
		Log: func(ctx *mserver.Context, err any) {
			fmt.Printf("panic 路径: %s", ctx.GetRequest().URL.String())
		},
	}
	c := mserver.NewCore()

	c.Get("/test/panic", func(ctx *mserver.Context) error {
		panic("test panic")
		return nil
	})
	c.Use(builder.Build())
	c.Start(":8081")
}
