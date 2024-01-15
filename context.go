package mserver

import (
	"net/http"
)

type Context struct {
	req  *http.Request
	resp http.ResponseWriter

	respStatusCode int
	params         map[string]string // url路由匹配的参数
	MatchedRoute   string

	index        int // 当前请求调用到调用链的哪个节点
	handlers     []HandleFunc
	CustomValues map[string]any // 自定义数据
}

// NewContext 初始化一个Context
func NewContext(r *http.Request, w http.ResponseWriter) *Context {
	return &Context{
		req:   r,
		resp:  w,
		index: -1,
	}
}

// SetHandlers 为context设置handlers
func (ctx *Context) SetHandlers(handlers []HandleFunc) {
	if ctx.handlers == nil {
		ctx.handlers = handlers
		return
	}
	ctx.handlers = append(ctx.handlers, handlers...)
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
