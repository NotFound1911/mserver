package test

import (
	"errors"
	"github.com/NotFound1911/mserver"
	"github.com/NotFound1911/mserver/session"
	"github.com/NotFound1911/mserver/session/cookie"
	"github.com/NotFound1911/mserver/session/memory"
	"github.com/google/uuid"
	"net/http"
	"testing"
	"time"
)

func TestManager(t *testing.T) {
	core := mserver.NewCore()
	m := session.Manager{
		CtxKey: "_sess",
		Store:  memory.NewStore(30 * time.Second),
		Adapter: cookie.NewAdapter("sessid",
			cookie.WithCookieOption(func(c *http.Cookie) {
				c.HttpOnly = true
			})),
	}
	core.Get("/login", func(ctx *mserver.Context) error {
		id := uuid.New()
		sess, err := m.InitSession(ctx, id.String())
		if err != nil {
			if errors.Is(err, http.ErrNoCookie) {
				ctx.SetStatus(http.StatusBadRequest).Json("cookie 不存在")
				return err
			}
			ctx.SetStatus(http.StatusInternalServerError).Json("unknown error")
			return err
		}
		// set value
		if err = sess.Set(ctx.GetRequest().Context(), "test_key", "test_value"); err != nil {
			ctx.SetStatus(http.StatusInternalServerError).Json("set session error")
			return err
		}
		ctx.SetStatus(http.StatusOK).Json("login successful")
		return nil
	})
	core.Get("/resource", func(ctx *mserver.Context) error {
		sess, err := m.GetSession(ctx)
		if err != nil {
			ctx.SetStatus(http.StatusInternalServerError)
			return err
		}
		val, err := sess.Get(ctx.GetRequest().Context(), "test_key")
		if err != nil {
			ctx.SetStatus(http.StatusInternalServerError)
			return err
		}
		ctx.SetStatus(http.StatusOK).Json(val)
		return nil
	})
	core.Get("/logout", func(ctx *mserver.Context) error {
		_ = m.DeleteSession(ctx)
		ctx.SetStatus(http.StatusOK).Text("delete session successful...")
		return nil
	})
	core.Use(func(next mserver.HandleFunc) mserver.HandleFunc {
		return func(ctx *mserver.Context) error {
			if ctx.GetRequest().URL.Path != "/login" {
				sess, err := m.GetSession(ctx)
				if err != nil {
					ctx.SetStatus(http.StatusUnauthorized).Json("get session error")
					return err
				}
				ctx.CustomValues["sess"] = sess
				_ = m.Update(ctx.GetRequest().Context(), sess.ID())
			}
			next(ctx)
			return nil
		}
	})
	core.Start(":8888")
}
