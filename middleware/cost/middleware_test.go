//go:build e2e

package cost

import (
	"github.com/NotFound1911/mserver"
	"net/http"
	"testing"
)

func TestMiddlewareBuilder_Build(t *testing.T) {
	c := mserver.NewCore()
	c.Get("/test/cost", func(ctx *mserver.Context) error {
		ctx.SetStatus(http.StatusOK).Json("test cost")
		return nil
	})
	c.Use(MiddlewareBuilder{}.Build())
	c.Start(":8081")
}
