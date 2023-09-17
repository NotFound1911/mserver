package session

import (
	"github.com/NotFound1911/mserver"
)

type Manager struct {
	Store
	Adapter
	CtxKey string
}

// GetSession 从 ctx 中拿到 Session
func (m *Manager) GetSession(ctx *mserver.Context) (Session, error) {
	if ctx.CustomValues == nil {
		ctx.CustomValues = map[string]any{}
	}
	val, ok := ctx.CustomValues[m.CtxKey]
	if ok {
		return val.(Session), nil
	}
	id, err := m.Extract(ctx.GetRequest())
	if err != nil {
		return nil, err
	}
	sess, err := m.Get(ctx.GetRequest().Context(), id)
	ctx.CustomValues[m.CtxKey] = sess
	return sess, nil
}

// InitSession 初始化session并注入到http中
func (m *Manager) InitSession(ctx *mserver.Context, id string) (Session, error) {
	sess, err := m.Create(ctx.GetRequest().Context(), id)
	if err != nil {
		return nil, err
	}
	if err = m.Inject(id, ctx.GetResponse()); err != nil {
		return nil, err
	}
	return sess, nil
}

// UpdateSession 更新session
func (m *Manager) UpdateSession(ctx *mserver.Context) (Session, error) {
	sess, err := m.GetSession(ctx)
	if err != nil {
		return nil, err
	}
	// 更新session过期时间
	err = m.Update(ctx.GetRequest().Context(), sess.ID())
	if err != nil {
		return nil, err
	}
	if err = m.Inject(sess.ID(), ctx.GetResponse()); err != nil {
		return nil, err
	}
	return sess, nil
}

func (m *Manager) DeleteSession(ctx *mserver.Context) error {
	sess, err := m.GetSession(ctx)
	if err != nil {
		return err
	}
	err = m.Store.Delete(ctx.GetRequest().Context(), sess.ID())
	if err != nil {
		return nil
	}
	return m.Adapter.Delete(ctx.GetResponse())
}
