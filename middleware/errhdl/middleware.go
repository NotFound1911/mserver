package errhdl

import (
	"github.com/NotFound1911/mserver"
)

type MiddlewareBuilder struct {
	// 这种设计只能返回固定的值
	// 不能做到动态渲染
	resp map[int][]byte
}

func NewMiddlewareBuilder() *MiddlewareBuilder {
	return &MiddlewareBuilder{
		resp: map[int][]byte{},
	}
}

func (m *MiddlewareBuilder) AddCode(status int, data []byte) *MiddlewareBuilder {
	m.resp[status] = data
	return m
}

func (m MiddlewareBuilder) Build() mserver.Middleware {
	return func(next mserver.HandleFunc) mserver.HandleFunc {
		return func(ctx *mserver.Context) error {
			next(ctx)
			resp, ok := m.resp[ctx.GetRespStatusCode()]
			if ok {
				// 篡改结果
				ctx.SetRespData(resp)
			}
			return nil
		}
	}
}
