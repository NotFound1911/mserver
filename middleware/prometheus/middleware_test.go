//go:build e2e

package prometheus

import (
	"github.com/NotFound1911/mserver"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"math/rand"
	"net/http"
	"testing"
	"time"
)

func TestMiddlewareBuilder_Build(t *testing.T) {
	build := MiddlewareBuilder{
		Name:      "test",
		Subsystem: "server",
		Namespace: "response",
	}
	core := mserver.NewCore()
	core.Get("/test/promethus", func(ctx *mserver.Context) error {
		val := rand.Intn(1000) + 1
		time.Sleep(time.Duration(val) + time.Millisecond)
		ctx.SetStatus(http.StatusOK).Json(User{
			Name: "test",
		})
		return nil
	})
	core.Use(build.Build())
	go func() {
		http.Handle("/metrics", promhttp.Handler())
		http.ListenAndServe(":8082", nil)
	}()
	core.Start(":8081")
}

type User struct {
	Name string
}
