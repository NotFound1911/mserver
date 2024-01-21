//go:build e2e

package recovery

import (
	"github.com/NotFound1911/mserver"
	"testing"
)

func TestMiddlewareBuilder_Build(t *testing.T) {
	c := mserver.NewCore()

	c.Get("/test/panic", func(ctx *mserver.Context) error {
		panic("test panic")
		return nil
	})
	c.Use(MiddlewareBuilder{}.Build())
	c.Start(":8081")
}
