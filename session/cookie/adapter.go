package cookie

import (
	"github.com/NotFound1911/mserver/session"
	"net/http"
)

type Option func(adapter *Adapter)

func WithCookieOption(opt func(c *http.Cookie)) Option {
	return func(adapter *Adapter) {
		adapter.cookieOpt = opt
	}
}

type Adapter struct {
	cookieName string
	cookieOpt  func(c *http.Cookie)
}

var _ session.Adapter = &Adapter{}

func NewAdapter(cookieName string, opts ...Option) *Adapter {
	adapter := &Adapter{
		cookieName: cookieName,
		cookieOpt: func(c *http.Cookie) {
		},
	}
	for _, opt := range opts {
		opt(adapter)
	}
	return adapter
}
func (a *Adapter) Inject(id string, writer http.ResponseWriter) error {
	cookie := &http.Cookie{
		Name:  a.cookieName,
		Value: id,
	}
	a.cookieOpt(cookie)
	http.SetCookie(writer, cookie)
	return nil
}
func (a *Adapter) Extract(req *http.Request) (string, error) {
	cookie, err := req.Cookie(a.cookieName)
	if err != nil {
		return "", err
	}
	return cookie.Value, err
}
func (a *Adapter) Delete(writer http.ResponseWriter) error {
	cookie := &http.Cookie{
		Name:   a.cookieName,
		MaxAge: -1,
	}
	a.cookieOpt(cookie)
	http.SetCookie(writer, cookie)
	return nil
}
