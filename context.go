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
func (c *Context) SetHandlers(handlers []HandleFunc) {
	if c.handlers == nil {
		c.handlers = handlers
		return
	}
	c.handlers = append(c.handlers, handlers...)
}

// 设置参数
func (c *Context) SetParams(params map[string]string) {
	c.params = params
}

// 核心函数，调用context的下一个函数
func (c *Context) Next() error {
	c.index++
	if c.index < len(c.handlers) {
		if err := c.handlers[c.index](c); err != nil {
			return err
		}
	}
	return nil
}

func (c *Context) GetRequest() *http.Request {
	return c.req
}
