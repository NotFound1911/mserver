package errhdl

import (
	"fmt"
	"github.com/NotFound1911/mserver"
	"net/http"
	"testing"
)

func TestMiddlewareBuilder_Build(t *testing.T) {
	c := mserver.NewCore()

	c.Get("/test/panic", func(ctx *mserver.Context) error {
		//panic("test panic")
		fmt.Println("666666666666")
		ctx.SetStatus(http.StatusOK).Json("66666")
		return nil
	})

	builder := NewMiddlewareBuilder()
	builder.AddCode(http.StatusNotFound, []byte(`
<html>
	<body>
		<h1>404:没有找到</h1>
	</body>
</html>
`)).
		AddCode(http.StatusBadRequest, []byte(`
<html>
	<body>
		<h1>500:请求错误</h1>
	</body>
</html>
`))

	c.Use(builder.Build())
	c.Start(":8081")
}
