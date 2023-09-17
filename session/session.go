package session

import (
	"context"
	"net/http"
)

// Session 存储和查找用户设置的数据
type Session interface {
	Get(ctx context.Context, key string) (string, error)
	Set(ctx context.Context, key string, val string) error
	ID() string
}

// Store 管理session的存储
type Store interface {
	Create(ctx context.Context, id string) (Session, error) // 创建session
	Update(ctx context.Context, id string) error            // 刷新session
	Delete(ctx context.Context, id string) error            // 删除session
	Get(ctx context.Context, id string) (Session, error)    // 查找session
}

// Adapter 适配层，不同的实现允许将 session id 存储在不同的地方
type Adapter interface {
	Inject(id string, writer http.ResponseWriter) error // 幂等操作, 将session id注入
	Extract(req *http.Request) (string, error)          // 将session id 从http.Request 进行提取
	Delete(writer http.ResponseWriter) error            // 将 session id从 http.ResponseWriter 中移除
}
