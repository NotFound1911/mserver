package mserver

import (
	"fmt"
	"net/http"
)

type Server interface {
	http.Handler
	// Start 启动服务器
	// addr 是监听地址。如果只指定端口，可以使用 ":8081"
	// 或者 "localhost:8082"
	Start(addr string) error

	// addRoute 注册一个路由
	// method 是 HTTP 方法
	addRoute(method string, path string, handler HandleFunc, mws ...Middleware)
}

var _ Server = &Core{}

type Core struct {
	router
	middlewares []Middleware
}

func NewCore() *Core {
	return &Core{
		router:      newRouter(),
		middlewares: []Middleware{},
	}
}

// ServeHTTP 处理请求的入口
func (c *Core) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	ctx := NewContext(request, writer)
	root := c.serve
	for i := len(c.middlewares) - 1; i >= 0; i-- {
		root = c.middlewares[i](root)
	}
	root(ctx)
}

func (c *Core) serve(ctx *Context) error {
	// 寻找路由
	mn, ok := c.FindRouteNodeByRequest(ctx.GetRequest())
	if !ok {
		// 未到路由 直接返回
		ctx.SetStatus(http.StatusNotFound).Text("%s not found", ctx.GetRequest().URL)
		return fmt.Errorf("未找到路由")
	}
	ctx.SetParams(mn.pathParams)
	ctx.MatchedRoute = mn.n.path
	handler := mn.n.handler
	for i := len(mn.matchMiddlewares) - 1; i >= 0; i-- {
		handler = mn.matchMiddlewares[i](handler)
	}
	handler(ctx)
	return nil
}

// Start 启动服务
func (c *Core) Start(addr string) error {
	return http.ListenAndServe(addr, c)
}

func (c *Core) StartTLS(addr, certFile, keyFile string) error {
	return http.ListenAndServeTLS(addr, certFile, keyFile, c)
}

func (c *Core) addRoute(method string, path string, handler HandleFunc, mws ...Middleware) {
	if err := c.router.addRoute(method, path, handler, mws...); err != nil {
		panic(fmt.Sprintf("add route err:%v", err))
	}
}

// 匹配路由
func (c *Core) FindRouteNodeByRequest(request *http.Request) (*matchNode, bool) {
	path := request.URL.Path
	method := request.Method
	return c.router.findRoute(method, path)
}

func (c *Core) Post(path string, handler HandleFunc) {
	c.addRoute(http.MethodPost, path, handler)
}
func (c *Core) Get(path string, handler HandleFunc) {
	c.addRoute(http.MethodGet, path, handler)
}
func (c *Core) Put(path string, handler HandleFunc) {
	c.addRoute(http.MethodPut, path, handler)
}
func (c *Core) Delete(path string, handler HandleFunc) {
	c.addRoute(http.MethodDelete, path, handler)
}

// 注册中间件
func (c *Core) Use(middlewares ...Middleware) {
	if c.middlewares == nil {
		c.middlewares = middlewares
		return
	}
	c.middlewares = append(c.middlewares, middlewares...)
}

func (c *Core) UsePath(method string, path string, mws ...Middleware) {
	c.addRoute(method, path, nil, mws...)
}
func (c *Core) Group(prefix string) Grouper {
	return NewGroup(c, prefix)
}
