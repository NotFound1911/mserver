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
	// 寻找路由
	mn, ok := c.FindRouteNodeByRequest(request)
	if !ok {
		// 未到路由 直接返回
		ctx.respStatusCode = 404
		return
	}

	allHandlers := make([]HandleFunc, 0, 8)
	for _, middleware := range c.middlewares {
		allHandlers = append(allHandlers, (HandleFunc)(middleware))
	}
	for _, middleware := range mn.matchMiddlewares {
		allHandlers = append(allHandlers, (HandleFunc)(middleware))
	}
	allHandlers = append(allHandlers, mn.n.handler)
	ctx.SetHandlers(allHandlers)
	// 设置路由参数
	ctx.SetParams(mn.pathParams)
	// 调用路由函数，如果返回err 代表存在内部错误，返回500状态码
	if err := ctx.Next(); err != nil {
		ctx.respStatusCode = 500
		return
	}
}

// 启动服务
func (c *Core) Start(addr string) error {
	return http.ListenAndServe(addr, c)
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
