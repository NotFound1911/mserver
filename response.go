package mserver

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"html/template"
	"net/http"
	"net/url"
)

type Responser interface {
	// Json Json输出
	Json(obj interface{}) Responser
	// Xml Xml输出
	Xml(obj interface{}) Responser
	// Html Html输出
	Html(template string, obj interface{}) Responser
	// Text string输出
	Text(format string, values ...interface{}) Responser
	// Redirect 重定向
	Redirect(path string) Responser
	// SetCookie 设置Cookie
	SetCookie(key string, val string, maxAge int, path, domain string, secure, httpOnly bool) Responser
	// SetStatus 设置状态码
	SetStatus(code int) Responser
	// SetHeader 设置header
	SetHeader(key string, val string) Responser
	// SetOkStatus 设置200状态
	SetOkStatus() Responser
}

func (ctx *Context) SetHeader(key string, val string) Responser {
	ctx.resp.Header().Add(key, val)
	return ctx
}
func (ctx *Context) Json(obj interface{}) Responser {
	byt, err := json.Marshal(obj)
	if err != nil {
		return ctx.SetStatus(http.StatusInternalServerError)
	}
	ctx.SetHeader("Content-Type", "application/json")
	ctx.respData = byt
	return ctx
}

func (ctx *Context) SetStatus(code int) Responser {
	ctx.respStatusCode = code
	return ctx
}
func (ctx *Context) SetOkStatus() Responser {
	ctx.resp.WriteHeader(http.StatusOK)
	return ctx
}
func (ctx *Context) Redirect(path string) Responser {
	http.Redirect(ctx.resp, ctx.req, path, http.StatusMovedPermanently)
	return ctx
}
func (ctx *Context) Xml(obj interface{}) Responser {
	byt, err := xml.Marshal(obj)
	if err != nil {
		return ctx.SetStatus(http.StatusInternalServerError)
	}
	ctx.SetHeader("Content-Type", "application/html")
	ctx.respData = byt
	return ctx
}
func (ctx *Context) Html(file string, obj interface{}) Responser {
	// 读取模版文件，创建template实例
	t, err := template.New("output").ParseFiles(file)
	if err != nil {
		return ctx
	}
	// 执行Execute方法将obj和模版进行结合
	if err := t.Execute(ctx.resp, obj); err != nil {
		return ctx
	}

	ctx.SetHeader("Content-Type", "application/html")
	return ctx
}
func (ctx *Context) Text(format string, values ...interface{}) Responser {
	out := fmt.Sprintf(format, values...)
	ctx.SetHeader("Content-Type", "application/text")
	ctx.respData = []byte(out)
	return ctx
}
func (ctx *Context) SetCookie(key string, val string, maxAge int, path string,
	domain string, secure bool, httpOnly bool) Responser {
	if path == "" {
		path = "/"
	}
	http.SetCookie(ctx.resp, &http.Cookie{
		Name:     key,
		Value:    url.QueryEscape(val),
		MaxAge:   maxAge,
		Path:     path,
		Domain:   domain,
		SameSite: 1,
		Secure:   secure,
		HttpOnly: httpOnly,
	})
	return ctx
}
