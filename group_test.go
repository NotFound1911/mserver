package mserver

import (
	"github.com/stretchr/testify/assert"
	"net/http"
	"reflect"
	"testing"
)

func Test_Group_route(t *testing.T) {
	core := NewCore()
	mockHandler1 := func(ctx *Context) error { return nil }
	mockHandler2 := func(ctx *Context) error { return nil }
	aGroup := core.Group("/a")
	{
		aGroup.Get("/:id", mockHandler2)
		aGroup.Put("/:id", mockHandler1)
		aGroup.Get("/test/list", mockHandler2)
		aGroup.Get("/test/*", mockHandler1)
		bGroup := aGroup.Group("/b")
		{
			bGroup.Get("/name", mockHandler2)
		}
	}
	testCases := []struct {
		name   string
		method string
		path   string
		found  bool
		mi     *matchNode
	}{
		{
			name:   "/a",
			method: http.MethodGet,
			path:   "/a",
			found:  false,
		},
		{
			name:   "/a/12",
			method: http.MethodGet,
			path:   "/a/12",
			found:  true,
			mi: &matchNode{
				n: &node{
					segment: ":id",
					handler: mockHandler2,
				},
				pathParams: map[string]string{"id": "12"},
			},
		},
		{
			name:   "/a/12",
			method: http.MethodPut,
			path:   "/a/12",
			found:  true,
			mi: &matchNode{
				n: &node{
					segment: ":id",
					handler: mockHandler1,
				},
				pathParams: map[string]string{"id": "12"},
			},
		},
		{
			name:   "/a/test/list",
			method: http.MethodGet,
			path:   "/a/test/list",
			found:  true,
			mi: &matchNode{
				n: &node{
					segment: "list",
					handler: mockHandler2,
				},
			},
		},
		{
			name:   "/a/test/all",
			method: http.MethodGet,
			path:   "/a/test/all",
			found:  true,
			mi: &matchNode{
				n: &node{
					segment: "all",
					handler: mockHandler1,
				},
			},
		},
		{
			name:   "/a/b/c",
			method: http.MethodGet,
			path:   "/a/b/c",
			found:  false,
		},
		{
			name:   "/a/b/name",
			method: http.MethodGet,
			path:   "/a/b/name",
			found:  true,
			mi: &matchNode{
				n: &node{
					segment: "name",
					handler: mockHandler2,
				},
			},
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mi, found := core.findRoute(tc.method, tc.path)
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
