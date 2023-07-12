//go:build e2e

package errhdl

import (
	"frame/web"
	"net/http"
	"testing"
)

func TestMiddlewareBuilder_Build(t *testing.T) {
	builder := NewMiddlewareBuilder()
	builder.AddCode(http.StatusNotFound, []byte(`
<html>
	<body>
		<h1> 页面只要不到阿, 兄弟儿 </h1>
	</body>
</html>

`)).
		AddCode(http.StatusBadRequest, []byte(`
<html>
	<body>
		<h1> 页面请求的不对阿, 兄弟儿 </h1>
	</body>
</html>
`))
	server := web.NewHTTPServer(web.ServerWithMiddleware(builder.Build()))
	server.Start(":8081")
}
