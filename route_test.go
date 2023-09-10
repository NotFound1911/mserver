package mserver

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"net/http"
	"reflect"
	"testing"
)

func Test_addRoute(t *testing.T) {
	testRoutes := []struct {
		method string
		path   string
	}{
		{
			method: http.MethodGet,
			path:   "/",
		},
		{
			method: http.MethodGet,
			path:   "/user",
		},
		{
			method: http.MethodGet,
			path:   "/user/home",
		},
		{
			method: http.MethodGet,
			path:   "/order/detail",
		},
		{
			method: http.MethodPost,
			path:   "/order/create",
		},
		{
			method: http.MethodPost,
			path:   "/login",
		},
		// 通配符测试用例
		{
			method: http.MethodGet,
			path:   "/order/*",
		},
		{
			method: http.MethodGet,
			path:   "/*",
		},
		{
			method: http.MethodGet,
			path:   "/*/*",
		},
		{
			method: http.MethodGet,
			path:   "/*/abc",
		},
		{
			method: http.MethodGet,
			path:   "/*/abc/*",
		},
		// 参数路由
		{
			method: http.MethodGet,
			path:   "/param/:id",
		},
		{
			method: http.MethodGet,
			path:   "/param/:id/detail",
		},
		{
			method: http.MethodGet,
			path:   "/param/:id/*",
		},
	}
	mockHandler := func(ctx *Context) error { return nil }
	r := newRouter()
	for _, tr := range testRoutes {
		r.addRoute(tr.method, tr.path, mockHandler)
	}
	wantRouter := &router{
		trees: map[string]*Tree{
			http.MethodGet: {
				Root: &node{
					path:    "/",
					segment: "",
					children: map[string]*node{
						"user": {
							path:    "/user",
							segment: "user",
							children: map[string]*node{
								"home": {
									path:    "/user/home",
									segment: "home",
									handler: mockHandler,
								},
							},
							handler: mockHandler,
						},
						"order": {
							path:    "/order",
							segment: "order",
							children: map[string]*node{
								"detail": {
									path:    "/order/detail",
									segment: "detail",
									handler: mockHandler,
								},
							},
							starChild: &node{
								path:    "/order/*",
								segment: "*",
								handler: mockHandler,
							},
						},
						"param": {
							path:    "/param",
							segment: "param",
							paramChild: &node{
								path:    "/param/:id",
								segment: ":id",
								starChild: &node{
									path:    "/param/*",
									segment: "*",
									handler: mockHandler,
								},
								children: map[string]*node{
									"detail": {
										path:    "/param/detail",
										segment: "detail",
										handler: mockHandler,
									},
								},
							},
						},
					},
					starChild: &node{
						path:    "/*",
						segment: "*",
						children: map[string]*node{
							"abc": {
								segment: "abc",
								path:    "/*/abc",
								starChild: &node{
									path:    "/*/abc/*",
									segment: "*",
									handler: mockHandler,
								},
								handler: mockHandler,
							},
						},
						starChild: &node{
							path:    "/*/*",
							segment: "*",
							handler: mockHandler,
						},
						handler: mockHandler,
					},
					handler: mockHandler,
				},
			},
			http.MethodPost: {
				Root: &node{
					path:    "/",
					segment: "",
					children: map[string]*node{
						"order": {
							path:    "/order",
							segment: "order",
							children: map[string]*node{
								"create": {
									path:    "/path/create",
									segment: "create",
									handler: mockHandler,
								},
							},
						},
						"login": {
							path:    "/login",
							segment: "login",
							handler: mockHandler,
						},
					},
				},
			},
		},
	}
	msg, ok := wantRouter.equal(r)
	assert.True(t, ok, msg)

	// 非法用例
	r = newRouter()
	// 空字符串
	{
		err := r.addRoute(http.MethodGet, "", mockHandler)
		assert.Errorf(t, err, "路由不允许为空")
	}
	// 前导没有 /
	{
		err := r.addRoute(http.MethodGet, "a/b/c", mockHandler)
		assert.Errorf(t, err, "路由必须以 / 开头")
	}
	// 后缀有 /
	{
		err := r.addRoute(http.MethodGet, "/a/b/c/", mockHandler)
		assert.Errorf(t, err, "路由不能以 / 结尾")
	}
	// 根节点重复注册
	{
		err := r.addRoute(http.MethodGet, "/", mockHandler)
		if err != nil {
			t.Errorf("err:%v is not nil", err)
		}
		err = r.addRoute(http.MethodGet, "/", mockHandler)
		assert.Errorf(t, err, "路由 [/] 冲突")
	}
	// 普通节点重复注册
	{
		err := r.addRoute(http.MethodGet, "/a/b/c", mockHandler)
		if err != nil {
			t.Errorf("err:%v is not nil", err)
		}
		err = r.addRoute(http.MethodGet, "/a/b/c", mockHandler)
		assert.Errorf(t, err, "路由 [/a/b/c] 冲突")
	}
	// 多个 /
	{
		err := r.addRoute(http.MethodGet, "/a//b", mockHandler)
		assert.Errorf(t, err, "非法路由:[/a//b]。不允许使用 //a/b, /a//b 之类的路由")
		err = r.addRoute(http.MethodGet, "//a/b", mockHandler)
		assert.Errorf(t, err, "非法路由:[//a/b]。不允许使用 //a/b, /a//b 之类的路由")
	}
	// 同时注册通配符路由和参数路由
	{
		err := r.addRoute(http.MethodGet, "/a/*", mockHandler)
		if err != nil {
			t.Errorf("err:%v is not nil", err)
		}
		err = r.addRoute(http.MethodGet, "/a/:id", mockHandler)
		assert.Errorf(t, err, "非法路由:[:id]，已有路径参数路由。不允许同时注册通配符路由和参数路由")
	}
	{
		err := r.addRoute(http.MethodGet, "/a/b/:id", mockHandler)
		if err != nil {
			t.Errorf("err:%v is not nil", err)
		}
		err = r.addRoute(http.MethodGet, "/a/b/*", mockHandler)
		assert.Errorf(t, err, "非法路由:[*]，已有路径参数路由。不允许同时注册通配符路由和参数路由")
	}
	r = newRouter()
	{
		err := r.addRoute(http.MethodGet, "/*", mockHandler)
		if err != nil {
			t.Errorf("err:%v is not nil", err)
		}
		err = r.addRoute(http.MethodGet, "/:id", mockHandler)
		assert.Errorf(t, err, "非法路由:[:id]，已有路径参数路由。不允许同时注册通配符路由和参数路由")
	}
	r = newRouter()
	{
		err := r.addRoute(http.MethodGet, "/:id", mockHandler)
		if err != nil {
			t.Errorf("err:%v is not nil", err)
		}
		err = r.addRoute(http.MethodGet, "/*", mockHandler)
		assert.Errorf(t, err, "非法路由:[*]，已有路径参数路由。不允许同时注册通配符路由和参数路由")
	}
	// 参数冲突
	{
		err := r.addRoute(http.MethodGet, "/a/b/c/:id", mockHandler)
		if err != nil {
			t.Errorf("err:%v is not nil", err)
		}
		err = r.addRoute(http.MethodGet, "/a/b/c/:name", mockHandler)
		assert.Errorf(t, err, "路由冲突，参数路由冲突，已有 [:id]，新注册 [:name]")
	}
}

func (r router) equal(y router) (string, bool) {
	for k, v := range r.trees {
		yv, ok := y.trees[k]
		if !ok {
			return fmt.Sprintf("目标 router 里面没有方法 %s 的路由树", k), false
		}
		str, ok := v.Root.equal(yv.Root)
		if !ok {
			return k + "-" + str, ok
		}
	}
	return "", true
}
func (n *node) equal(y *node) (string, bool) {
	if y == nil {
		return "目标节点为 nil", false
	}
	if n.segment != y.segment {
		return fmt.Sprintf("[%s] 节点 segment 不相等 x [%s], y [%s]", n.segment, n.segment, y.segment), false
	}

	nhv := reflect.ValueOf(n.handler)
	yhv := reflect.ValueOf(y.handler)
	if nhv != yhv {
		return fmt.Sprintf("[%s] 节点 handler 不相等 x [%s], y [%s]", n.segment, nhv.Type().String(), yhv.Type().String()), false
	}

	if len(n.children) != len(y.children) {
		return fmt.Sprintf("%s 子节点长度不等", n.segment), false
	}
	if len(n.children) == 0 {
		return "", true
	}

	if n.starChild != nil {
		str, ok := n.starChild.equal(y.starChild)
		if !ok {
			return fmt.Sprintf("%s 通配符节点不匹配 %s", n.segment, str), false
		}
	}

	for k, v := range n.children {
		yv, ok := y.children[k]
		if !ok {
			return fmt.Sprintf("%s 目标节点缺少子节点 %s", n.segment, k), false
		}
		str, ok := v.equal(yv)
		if !ok {
			return n.segment + "-" + str, ok
		}
	}
	return "", true
}
func Test_findRoute(t *testing.T) {
	testRoutes := []struct {
		method string
		path   string
	}{
		{
			method: http.MethodGet,
			path:   "/",
		},
		{
			method: http.MethodGet,
			path:   "/user",
		},
		{
			method: http.MethodPost,
			path:   "/order/create",
		},
		{
			method: http.MethodGet,
			path:   "/user/*/home",
		},
		{
			method: http.MethodPost,
			path:   "/order/*",
		},
		// 参数路由
		{
			method: http.MethodGet,
			path:   "/param/:id",
		},
		{
			method: http.MethodGet,
			path:   "/param/:id/detail",
		},
		{
			method: http.MethodGet,
			path:   "/param/:id/*",
		},
	}

	mockHandler := func(ctx *Context) error { return nil }

	testCases := []struct {
		name   string
		method string
		path   string
		found  bool
		mi     *matchNode
	}{
		{
			name:   "method not found",
			method: http.MethodHead,
		},
		{
			name:   "path not found",
			method: http.MethodGet,
			path:   "/abc",
		},
		{
			name:   "root",
			method: http.MethodGet,
			path:   "/",
			found:  true,
			mi: &matchNode{
				n: &node{
					segment: "",
					handler: mockHandler,
				},
			},
		},
		{
			name:   "user",
			method: http.MethodGet,
			path:   "/user",
			found:  true,
			mi: &matchNode{
				n: &node{
					segment: "user",
					handler: mockHandler,
				},
			},
		},
		{
			name:   "no handler",
			method: http.MethodPost,
			path:   "/order",
			found:  false,
			mi: &matchNode{
				n: &node{
					segment: "order",
				},
			},
		},
		{
			name:   "two layer",
			method: http.MethodPost,
			path:   "/order/create",
			found:  true,
			mi: &matchNode{
				n: &node{
					segment: "create",
					handler: mockHandler,
				},
			},
		},
		// 通配符匹配
		{
			// 命中/order/*
			name:   "star match",
			method: http.MethodPost,
			path:   "/order/delete",
			found:  true,
			mi: &matchNode{
				n: &node{
					segment: "*",
					handler: mockHandler,
				},
			},
		},
		{
			// 命中通配符在中间的
			// /user/*/home
			name:   "star in middle",
			method: http.MethodGet,
			path:   "/user/Tom/home",
			found:  true,
			mi: &matchNode{
				n: &node{
					segment: "home",
					handler: mockHandler,
				},
			},
		},
		{
			// 比 /order/* 多了一段
			name:   "overflow",
			method: http.MethodPost,
			path:   "/order/delete/123",
		},
		// 参数匹配
		{
			// 命中 /param/:id
			name:   ":id",
			method: http.MethodGet,
			path:   "/param/123",
			found:  true,
			mi: &matchNode{
				n: &node{
					segment: ":id",
					handler: mockHandler,
				},
				pathParams: map[string]string{"id": "123"},
			},
		},
		{
			// 命中 /param/:id/*
			name:   ":id*",
			method: http.MethodGet,
			path:   "/param/123/abc",
			found:  true,
			mi: &matchNode{
				n: &node{
					segment: "*",
					handler: mockHandler,
				},
				pathParams: map[string]string{"id": "123"},
			},
		},

		{
			// 命中 /param/:id/detail
			name:   ":id*",
			method: http.MethodGet,
			path:   "/param/123/detail",
			found:  true,
			mi: &matchNode{
				n: &node{
					segment: "detail",
					handler: mockHandler,
				},
				pathParams: map[string]string{"id": "123"},
			},
		},
	}
	r := newRouter()
	for _, tr := range testRoutes {
		r.addRoute(tr.method, tr.path, mockHandler)
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mi, found := r.findRoute(tc.method, tc.path)
			assert.Equal(t, tc.found, found)
			if !found {
				return
			}
			assert.Equal(t, tc.mi.pathParams, mi.pathParams)
			n := mi.n
			wantVal := reflect.ValueOf(tc.mi.n.handler)
			nVal := reflect.ValueOf(n.handler)
			assert.Equal(t, wantVal, nVal)
		})
	}
}
func Test_findRoute_mws(t *testing.T) {
	mockHandler := func(ctx *Context) error { return nil }
	mockHandler1 := func(next HandleFunc) HandleFunc { return next }
	mockHandler2 := func(next HandleFunc) HandleFunc { return next }
	mockHandler3 := func(next HandleFunc) HandleFunc { return next }
	mockHandler4 := func(next HandleFunc) HandleFunc { return next }
	mockHandler5 := func(next HandleFunc) HandleFunc { return next }
	mockHandler6 := func(next HandleFunc) HandleFunc { return next }
	// 可路由中间件测试
	testRoutes := []struct {
		method string
		path   string
		mws    []Middleware
	}{
		{
			method: http.MethodGet,
			path:   "/a/b",
			mws:    []Middleware{mockHandler1},
		},
		{
			method: http.MethodGet,
			path:   "/a/b/c",
			mws:    []Middleware{mockHandler2},
		},
		{
			method: http.MethodGet,
			path:   "/a/*",
			mws:    []Middleware{mockHandler3},
		},
		{
			method: http.MethodGet,
			path:   "/a/b/*",
			mws:    []Middleware{mockHandler4},
		},
		{
			method: http.MethodGet,
			path:   "/a/*/d",
			mws:    []Middleware{mockHandler5, mockHandler6},
		},
		{
			method: http.MethodGet,
			path:   "/a1/:id",
			mws:    []Middleware{mockHandler1},
		},
		{
			method: http.MethodGet,
			path:   "/a1/123/c",
			mws:    []Middleware{mockHandler2},
		},
		{
			method: http.MethodGet,
			path:   "/a1/:id/c",
			mws:    []Middleware{mockHandler3},
		},
	}
	testCases := []struct {
		name   string
		method string
		path   string
		found  bool
		mi     *matchNode
	}{
		{
			name:   "/a/b 匹配",
			method: http.MethodGet,
			path:   "/a/b",
			found:  true,
			mi: &matchNode{
				n: &node{
					segment: "b",
					handler: mockHandler,
				},
				matchMiddlewares: []Middleware{mockHandler3, mockHandler1},
			},
		},
		{
			name:   "/a/b/c 匹配",
			method: http.MethodGet,
			path:   "/a/b/c",
			found:  true,
			mi: &matchNode{
				n: &node{
					segment: "c",
					handler: mockHandler,
				},
				matchMiddlewares: []Middleware{mockHandler3, mockHandler1, mockHandler4, mockHandler2},
			},
		},
		{
			name:   "/a/b/d 匹配",
			method: http.MethodGet,
			path:   "/a/b/d",
			found:  true,
			mi: &matchNode{
				n: &node{
					segment: "c",
					handler: mockHandler,
				},
				matchMiddlewares: []Middleware{mockHandler3, mockHandler1, mockHandler5, mockHandler6, mockHandler4},
			},
		},
		{
			name:   "/a/e 匹配",
			method: http.MethodGet,
			path:   "/a/e",
			found:  true,
			mi: &matchNode{
				n: &node{
					segment: "e",
					handler: mockHandler,
				},
				matchMiddlewares: []Middleware{mockHandler3},
			},
		},
		{
			name:   "/a1/:id 匹配",
			method: http.MethodGet,
			path:   "/a1/122",
			found:  true,
			mi: &matchNode{
				n: &node{
					segment: ":id",
					handler: mockHandler,
				},
				pathParams:       map[string]string{"id": "122"},
				matchMiddlewares: []Middleware{mockHandler1},
			},
		},
		{
			name:   "/a1/123/c 匹配",
			method: http.MethodGet,
			path:   "/a1/123/c",
			found:  true,
			mi: &matchNode{
				n: &node{
					segment: "c",
					handler: mockHandler,
				},
				matchMiddlewares: []Middleware{mockHandler1, mockHandler3, mockHandler2},
			},
		},
		{
			name:   "/a1/00/c 匹配",
			method: http.MethodGet,
			path:   "/a1/00/c",
			found:  true,
			mi: &matchNode{
				n: &node{
					segment: "c",
					handler: mockHandler,
				},
				pathParams:       map[string]string{"id": "00"},
				matchMiddlewares: []Middleware{mockHandler1, mockHandler3},
			},
		},
	}
	r := newRouter()
	for _, tr := range testRoutes {
		r.addRoute(tr.method, tr.path, mockHandler, tr.mws...)
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mi, found := r.findRoute(tc.method, tc.path)
			assert.Equal(t, tc.found, found)
			if !found {
				return
			}
			assert.Equal(t, tc.mi.pathParams, mi.pathParams)
			n := mi.n
			wantVal := reflect.ValueOf(tc.mi.n.handler)
			nVal := reflect.ValueOf(n.handler)
			assert.Equal(t, wantVal, nVal)
			assert.Equal(t, len(tc.mi.matchMiddlewares), len(mi.matchMiddlewares))
			for i := 0; i < len(tc.mi.matchMiddlewares); i++ {
				wantValMw := reflect.ValueOf(tc.mi.matchMiddlewares[i])
				nValMw := reflect.ValueOf(mi.matchMiddlewares[i])
				assert.Equal(t, wantValMw, nValMw)
			}
		})
	}
}
