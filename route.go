package mserver

import (
	"errors"
	"fmt"
	"strings"
)

type router struct {
	trees map[string]*Tree
}

func newRouter() router {
	return router{
		trees: map[string]*Tree{},
	}
}

// addRoute 限制:
// 1.不允许空路由
// 2.不运行重复注册
// 3.路由必须/开头
// 4.路由不能以/结尾
func (r *router) addRoute(method string, path string, handler HandleFunc, mws ...Middleware) error {
	if path == "" {
		return errors.New("路由不允许为空")
	}
	if path[0] != '/' {
		return errors.New("路由必须以 / 开头")
	}
	if path != "/" && path[len(path)-1] == '/' {
		return errors.New("路由不能以 / 结尾")
	}
	root, ok := r.trees[method]
	if !ok { // 全新的 http 方法
		root = &Tree{
			Root: &node{
				path: "/",
			},
		}
		r.trees[method] = root
	}
	if path == "/" {
		if root.Root.handler != nil {
			return errors.New("路由 [/] 冲突")
		}
		root.Root.handler = handler
		root.Root.mws = mws
		return nil
	}
	pNode := root.Root
	segments := strings.Split(path[1:], "/")
	for _, seg := range segments {
		if seg == "" {
			return fmt.Errorf("非法路由:[%s]。不允许使用 //a/b, /a//b 之类的路由", path)
		}
		// 查询子节点
		childNode, err := pNode.findeChildNode(seg) // 匹配子节点
		if err != nil {
			return err
		}
		if childNode == nil { // 空节点 需要创建
			cnode := newNode()
			childNode = cnode
			childNode.segment = seg
		}
		// 节点连接
		switch connectNodeWay(seg) {
		case 0:
			pNode.children[seg] = childNode
		case 1:
			pNode.paramChild = childNode
		case 2:
			pNode.starChild = childNode
		}
		pNode = childNode
	}
	if pNode.handler != nil {
		return fmt.Errorf("路由 [%s] 冲突", path)
	}
	pNode.handler = handler
	pNode.path = path
	pNode.mws = mws
	return nil
}

// 连接方式
// 0: children 普通
// 1: paramChild 参数
// 2: starChild 通配符路径
func connectNodeWay(seg string) int {
	if seg == "*" {
		return 2
	}
	if seg[0] == ':' {
		return 1
	}
	return 0
}

// 查询路由
// 返回的 node 内部 HandleFunc 不为 nil 才算是注册了路由
func (r *router) findRoute(method string, path string) (*matchNode, bool) {
	tree, ok := r.trees[method]
	if !ok {
		return nil, false
	}
	if path == "/" {
		return &matchNode{n: tree.Root, matchMiddlewares: tree.Root.mws}, tree.Root.handler != nil
	}
	fNode := tree.matchNode(path)
	if fNode == nil || fNode.n.handler == nil {
		return fNode, false
	}
	return fNode, true
}
