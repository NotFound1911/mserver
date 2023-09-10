package mserver

import (
	"fmt"
	"net/http"
	"testing"
)

func Test_Core_Server(t *testing.T) {
	core := NewCore()
	registerRouter(core)
	t.Logf("register route successful...")
	core.Start(":8888")
	for {

	}
}

// 注册路由规则
func registerRouter(core *Core) {
	mockHandler := func(ctx *Context) error {
		ctx.respStatusCode = 200
		ctx.resp.Write([]byte("this is mockHandler"))
		return nil
	}
	core.addRoute(http.MethodGet, "/", mockHandler)
	core.addRoute(http.MethodGet, "/user/test", mockHandler)
}

func Test_Core_Server_route(t *testing.T) {
	core := NewCore()
	mw := func(next HandleFunc) HandleFunc {
		return func(ctx *Context) error {
			fmt.Println("this is mid")
			if err := next(ctx); err != nil {
				return err
			}
			return nil
		}
	}
	core.Get("/user/home", func(ctx *Context) error {
		ctx.SetStatus(http.StatusOK).Text("this is /usr/home")
		return nil
	})
	core.Get("/user/school", func(ctx *Context) error {
		ctx.SetStatus(http.StatusOK).Text("this is /user/school")
		return nil
	})
	core.UsePath(http.MethodGet, "/user/*", mw)
	core.Start(":8888")
}
