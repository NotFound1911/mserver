package mserver

type HandleFunc func(ctx *Context) error
type Middleware func(handleFunc HandleFunc) HandleFunc
