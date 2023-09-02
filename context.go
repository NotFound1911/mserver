package mserver

import (
	"net/http"
)

type Context struct {
	req  *http.Request
	resp http.ResponseWriter

	respStatusCode int
	params         map[string]string // url路由匹配的参数

	index    int // 当前请求调用到调用链的哪个节点
	handlers []HandleFunc
}

// NewContext 初始化一个Context
func NewContext(r *http.Request, w http.ResponseWriter) *Context {
	return &Context{
		req:   r,
		resp:  w,
		index: -1,
	}
}

// 为context设置handlers
func (ctx *Context) SetHandlers(handlers []HandleFunc) {
	if ctx.handlers == nil {
		ctx.handlers = handlers
		return
	}
	ctx.handlers = append(ctx.handlers, handlers...)
}

// 设置参数
func (ctx *Context) SetParams(params map[string]string) {
	ctx.params = params
}

// 核心函数，调用context的下一个函数
func (ctx *Context) Next() error {
	ctx.index++
	if ctx.index < len(ctx.handlers) {
		if err := ctx.handlers[ctx.index](ctx); err != nil {
			return err
		}
	}
	return nil
}
