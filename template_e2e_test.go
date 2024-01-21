//go:build e2e

package mserver

import (
	"github.com/stretchr/testify/require"
	"html/template"
	"log"
	"testing"
)

func TestLoginPage(t *testing.T) {
	tpl, err := template.ParseGlob("internal/testdata/tpls/*.gohtml")
	require.NoError(t, err)
	engine := &GoTemplateEngine{
		T: tpl,
	}

	c := NewCore(CoreWithTemplateEngine(engine))
	c.Get("/login", func(ctx *Context) error {
		err := ctx.Render("login.gohtml", nil)
		if err != nil {
			log.Println(err)
		}
		return nil
	})
	c.Start(":8081")
}
