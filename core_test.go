package mserver

import (
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
