package mserver

import "testing"

func Test_findeChildNode(t *testing.T) {
	root := &node{
		segment: "/",
		handler: func(ctx *Context) error {
			return nil
		},
		children: map[string]*node{
			"foo": {
				path:    "/foo",
				segment: "foo",
				handler: func(ctx *Context) error {
					return nil
				},
				children: nil,
			},
			"FOO": {
				path:    "/FOO",
				segment: "FOO",
				handler: func(ctx *Context) error {
					return nil
				},
				children: nil,
				starChild: &node{
					path:    "/FOO/*",
					segment: "*",
					handler: func(ctx *Context) error {
						return nil
					},
					children: nil,
				},
			},
		},
		paramChild: &node{
			path:    "/:id",
			segment: ":id",
			handler: func(ctx *Context) error {
				return nil
			},
			children: nil,
		},
		//starChild: &node{
		//	path:    "/*",
		//	segment: "*",
		//	handler: func(ctx *Context) error {
		//		return nil
		//	},
		//	children: nil,
		//},
	}
	{
		n, err := root.findeChildNode("FOO")
		t.Logf("node path:%s", n.path)
		if err != nil {
			t.Errorf("err:%v", err)
		}
	}
	{
		n, err := root.findeChildNode("foo")
		t.Logf("node path:%s", n.path)
		if err != nil {
			t.Errorf("err:%v", err)
		}
	}
	{
		n, err := root.findeChildNode(":id")
		t.Logf("node path:%s", n.path)
		if err != nil {
			t.Errorf("err:%v", err)
		}
	}
}

func Test_matchNode(t *testing.T) {
	root := &node{
		children: map[string]*node{
			"FOO": {
				segment: "FOO",
				children: map[string]*node{
					"BAR": {
						segment:  "BAR",
						children: map[string]*node{},
					},
				},
			},
		},
		paramChild: &node{
			segment: ":id",
			handler: func(ctx *Context) error {
				return nil
			},
			children: nil,
			starChild: &node{
				segment: "/:id/*",
				handler: func(ctx *Context) error {
					return nil
				},
				children: nil,
			},
		},
	}
	tree := &Tree{Root: root}
	{
		node := tree.matchNode("FOO/BAR")
		if node == nil {
			t.Error("match normal node error")
		} else {
			t.Logf("node segment:%s", node.n.segment)
		}

	}
	{
		node := tree.matchNode("ttt")
		if node == nil {
			t.Error("match normal node error")
		} else {
			t.Logf("node segment:%s", node.n.segment)
		}
	}
	{
		node := tree.matchNode("ttt/t1")
		if node == nil {
			t.Error("match normal node error")
		} else {
			t.Logf("node segment:%s", node.n.segment)
		}
	}
}
