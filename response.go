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

func (c *Context) SetHeader(key string, val string) Responser {
	c.resp.Header().Add(key, val)
	return c
}
func (c *Context) Json(obj interface{}) Responser {
	byt, err := json.Marshal(obj)
	if err != nil {
		return c.SetStatus(http.StatusInternalServerError)
	}
	c.SetHeader("Content-Type", "application/json")
	c.resp.Write(byt)
	return c
}

func (c *Context) SetStatus(code int) Responser {
	c.resp.WriteHeader(code)
	return c
}
func (c *Context) SetOkStatus() Responser {
	c.resp.WriteHeader(http.StatusOK)
	return c
}
func (c *Context) Redirect(path string) Responser {
	http.Redirect(c.resp, c.req, path, http.StatusMovedPermanently)
	return c
}
func (c *Context) Xml(obj interface{}) Responser {
	byt, err := xml.Marshal(obj)
	if err != nil {
		return c.SetStatus(http.StatusInternalServerError)
	}
	c.SetHeader("Content-Type", "application/html")
	c.resp.Write(byt)
	return c
}
func (c *Context) Html(file string, obj interface{}) Responser {
	// 读取模版文件，创建template实例
	t, err := template.New("output").ParseFiles(file)
	if err != nil {
		return c
	}
	// 执行Execute方法将obj和模版进行结合
	if err := t.Execute(c.resp, obj); err != nil {
		return c
	}

	c.SetHeader("Content-Type", "application/html")
	return c
}
func (c *Context) Text(format string, values ...interface{}) Responser {
	out := fmt.Sprintf(format, values...)
	c.SetHeader("Content-Type", "application/text")
	c.resp.Write([]byte(out))
	return c
}
func (c *Context) SetCookie(key string, val string, maxAge int, path string,
	domain string, secure bool, httpOnly bool) Responser {
	if path == "" {
		path = "/"
	}
	http.SetCookie(c.resp, &http.Cookie{
		Name:     key,
		Value:    url.QueryEscape(val),
		MaxAge:   maxAge,
		Path:     path,
		Domain:   domain,
		SameSite: 1,
		Secure:   secure,
		HttpOnly: httpOnly,
	})
	return c
}
