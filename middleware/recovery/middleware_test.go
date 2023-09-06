package recovery

import (
	"github.com/NotFound1911/mserver"
	"testing"
)

func TestMiddlewares_Recover(t *testing.T) {
	c := mserver.NewCore()
	c.Get("/test/panic", func(ctx *mserver.Context) error {
		panic("test panic")
		return nil
	})
	c.Use(Recovery())
	c.Start(":8081")
}
