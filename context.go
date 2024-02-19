package mserver

import (
	"context"
	"net/http"
	"time"
)

type Context struct {
	req  *http.Request
	resp http.ResponseWriter

	// 注意为middleware读写使用
	respStatusCode int
	respData       []byte

	params       map[string]string // url路由匹配的参数
	MatchedRoute string

	CustomValues map[string]any // 自定义数据

	tplEngine TemplateEngine
	// ContextWithFallback enable fallback Context.Deadline(), Context.Done(), Context.Err() and Context.Value()
	// when Context.Request.Context() is not nil.
	ContextWithFallback bool
}

var _ context.Context = &Context{}

func (ctx *Context) Deadline() (deadline time.Time, ok bool) {
	if !ctx.ContextWithFallback || ctx.req == nil || ctx.req.Context() == nil {
		return
	}
	return ctx.req.Context().Deadline()
}

func (ctx *Context) Done() <-chan struct{} {
	if !ctx.ContextWithFallback || ctx.req == nil || ctx.req.Context() == nil {
		return nil
	}
	return ctx.req.Context().Done()
}

func (ctx *Context) Err() error {
	if !ctx.ContextWithFallback || ctx.req == nil || ctx.req.Context() == nil {
		return nil
	}
	return ctx.req.Context().Err()
}

func (ctx *Context) Value(key any) any {
	if !ctx.ContextWithFallback || ctx.req == nil || ctx.req.Context() == nil {
		return nil
	}
	return ctx.req.Context().Value(key)
}

func newContext() *Context {
	return &Context{}
}
func (ctx *Context) reset() {
	ctx.respStatusCode = 0
	ctx.respData = nil
	ctx.params = nil // url路由匹配的参数
	ctx.MatchedRoute = ""
	ctx.CustomValues = nil // 自定义数据
	ctx.tplEngine = nil
	ctx.ContextWithFallback = false
}

// Render 渲染
func (ctx *Context) Render(tplName string, data any) error {
	// 不要这样子去做

	var err error
	respData, err := ctx.tplEngine.Render(ctx.GetRequest().Context(), tplName, data)
	ctx.SetRespData(respData)
	if err != nil {
		ctx.SetStatus(http.StatusInternalServerError)
		return err
	}
	ctx.SetStatus(http.StatusOK)
	return nil
}

// SetParams 设置参数
func (ctx *Context) SetParams(params map[string]string) {
	ctx.params = params
}

func (ctx *Context) GetRequest() *http.Request {
	return ctx.req
}
func (ctx *Context) GetResponse() http.ResponseWriter {
	return ctx.resp
}

func (ctx *Context) SetRequest(req *http.Request) {
	ctx.req = req
}

func (ctx *Context) GetRespStatusCode() int {
	return ctx.respStatusCode
}
func (ctx *Context) SetRespData(data []byte) {
	ctx.respData = data
}
