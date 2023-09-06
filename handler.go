package mserver

type HandleFunc func(ctx *Context) error
type Middleware HandleFunc // 等价的 只是为了区分
