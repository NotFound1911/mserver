package cost

import (
	"github.com/NotFound1911/mserver"
	"testing"
)

func TestMiddleware_Cost(t *testing.T) {
	c := mserver.NewCore()
	c.Get("/test/cost", func(ctx *mserver.Context) error {
		return nil
	})
	c.Use(Cost())
	c.Start(":8081")
}
