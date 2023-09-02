package mserver

import (
	"fmt"
	"strings"
)

// Tree 代表树结构
type Tree struct {
	Root *node
}

// 匹配节点
func (t *Tree) matchNode(path string) *matchNode {
	segs := strings.Split(strings.Trim(path, "/"), "/")
	pNode := t.Root
	mn := &matchNode{}
	for _, seg := range segs {
		cNode, isParam := pNode.childOf(seg)
		if cNode == nil {
			return nil
		}
		if isParam {
			mn.addValue(cNode.segment[1:], seg)
		}
		pNode = cNode
	}
	mn.n = pNode
	mn.matchMiddlewares = t.Root.findMdls(segs)
	return mn
}
func (root *node) findMdls(segs []string) []Middleware {
	queue := []*node{root}
	res := make([]Middleware, 0, 16)
	for i := 0; i < len(segs); i++ {
		seg := segs[i]
		var children []*node
		for _, cur := range queue {
			if len(cur.mws) > 0 {
				res = append(res, cur.mws...)
			}
			children = append(children, cur.childrenOf(seg)...)
		}
		queue = children
	}

	for _, cur := range queue {
		if len(cur.mws) > 0 {
			res = append(res, cur.mws...)
		}
	}
	return res
}

func (n *node) childrenOf(path string) []*node {
	res := make([]*node, 0, 4)
	var static *node
	if n.children != nil {
		static = n.children[path]
	}
	if n.starChild != nil {
		res = append(res, n.starChild)
	}
	if n.paramChild != nil {
		res = append(res, n.paramChild)
	}
	if static != nil {
		res = append(res, static)
	}
	return res
}

// 代表路由树的节点
type node struct {
	path    string
	segment string // uri中的字符串
	// handler 命中路由之后执行的逻辑
	handler HandleFunc
	// 注册在该节点上的 middleware
	mws []Middleware
	// 该节点匹配到的middleware 即实际运行的
	matchedMdls []Middleware
	// children 子节点
	// 子节点的 path => node
	children map[string]*node

	// 通配符 * 表达的节点，任意匹配
	starChild *node

	paramChild *node
}

func newNode() *node {
	return &node{
		segment:    "",
		children:   map[string]*node{},
		handler:    nil,
		starChild:  nil,
		paramChild: nil,
	}
}

// 过滤下一层满足segment规则的子节点.
// 首先会判断 path 是不是通配符路径p[]\
// 其次判断 path 是不是参数路径，即以 : 开头的路径
// 最后0会从 children 里面查找，
func (n *node) findeChildNode(seg string) (*node, error) {
	if seg == "*" {
		if n.paramChild != nil {
			return nil, fmt.Errorf("非法路由:[%s]，已有路径参数路由。不允许同时注册通配符路由和参数路由", seg)
		}
		return n.starChild, nil
	}
	// 以 : 开头，我们认为是参数路由
	if seg[0] == ':' {
		if n.starChild != nil {
			return nil, fmt.Errorf("非法路由:[%s]，已有路径参数路由。不允许同时注册通配符路由和参数路由", seg)
		}
		if n.paramChild != nil {
			if n.paramChild.segment != seg {
				return nil, fmt.Errorf("路由冲突，参数路由冲突，已有 [%s]，新注册 [%s]", n.paramChild.segment, seg)
			}
		} else {
			n.paramChild = &node{segment: seg}
		}
		return n.paramChild, nil
	}
	if n.children == nil {
		n.children = map[string]*node{}
	}
	return n.children[seg], nil
}

// child 返回子节点
// 匹配优先级: 1. 静态路由 2.参数路由 3.通配符
//
//	返回值：匹配的子节点 是否命中参数路由
func (n *node) childOf(path string) (*node, bool) {
	if n.children == nil {
		if n.paramChild != nil {
			return n.paramChild, true
		}
		return n.starChild, false
	}
	res, ok := n.children[path]
	if !ok {
		if n.paramChild != nil {
			return n.paramChild, true
		}
		return n.starChild, false
	}
	return res, false
}

type matchNode struct {
	n                *node
	pathParams       map[string]string
	matchMiddlewares []Middleware
}

func (m *matchNode) addValue(key string, value string) {
	if m.pathParams == nil {
		// 大多数情况，参数路径只会有一段
		m.pathParams = map[string]string{key: value}
	}
	m.pathParams[key] = value
}
