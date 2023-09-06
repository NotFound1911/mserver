package mserver

// Grouper 前缀分组
// 实现HttpMethod方法
type Grouper interface {
	Post(path string, handler HandleFunc)
	Put(path string, handler HandleFunc)
	Get(path string, handler HandleFunc)
	Delete(path string, handler HandleFunc)
	// Group 实现嵌套
	Group(string string) Grouper
}

var _ Grouper = &Group{}

type Group struct {
	core   *Core
	prefix string // 这个group的通用前缀
	parent *Group // 指向上一个ServerGroup
}

func NewGroup(core *Core, prefix string) *Group {
	return &Group{
		core:   core,
		parent: nil,
		prefix: prefix,
	}
}

func (g *Group) Post(path string, handler HandleFunc) {
	path = g.getAbsolutePrefix() + path
	g.core.Post(path, handler)
}
func (g *Group) Put(path string, handler HandleFunc) {
	path = g.getAbsolutePrefix() + path
	g.core.Put(path, handler)
}
func (g *Group) Get(path string, handler HandleFunc) {
	path = g.getAbsolutePrefix() + path
	g.core.Get(path, handler)
}
func (g *Group) Delete(path string, handler HandleFunc) {
	path = g.getAbsolutePrefix() + path
	g.core.Delete(path, handler)
}

func (g *Group) Group(uri string) Grouper {
	childGroup := NewGroup(g.core, uri)
	childGroup.parent = g
	return childGroup
}

// 获取当前group的绝对路径
func (g *Group) getAbsolutePrefix() string {
	if g.parent == nil {
		return g.prefix
	}
	return g.parent.getAbsolutePrefix() + g.prefix
}

func (g *Group) getCore() *Core {
	return g.core
}
