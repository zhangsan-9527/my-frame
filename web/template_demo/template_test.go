package template_demo

import (
	"bytes"
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"html/template"
	"testing"
)

func TestHelloWorld(t *testing.T) {
	type User struct {
		Name string
	}

	tpl := template.New("hello-world")
	// 使用.来访问数据，代表的是当前作用域的当前对象，类似于 Java 的 this、Python 的 self.
	// 所以.Name 代表的是访问传入的 User 的 Name 字段
	parse, err := tpl.Parse(`Hello, {{ .Name }}`)
	require.NoError(t, err)
	buffer := &bytes.Buffer{}
	// 指针结构体都行
	err = parse.Execute(buffer, &User{Name: "Tom"})
	require.NoError(t, err)
	assert.Equal(t, `Hello, Tom`, buffer.String())
}

func TestMapData(t *testing.T) {

	tpl := template.New("hello-world")
	// 使用.来访问数据，代表的是当前作用域的当前对象，类似于 Java 的 this、Python 的 self.
	// 所以.Name 代表的是访问传入的 User 的 Name 字段
	parse, err := tpl.Parse(`Hello, {{ .Name }}`)
	require.NoError(t, err)
	buffer := &bytes.Buffer{}
	// 指针结构体都行
	err = parse.Execute(buffer, map[string]string{"Name": "zhangsan"})
	require.NoError(t, err)
	assert.Equal(t, `Hello, zhangsan`, buffer.String())
}

func TestSlice(t *testing.T) {
	tpl := template.New("hello-world")
	// 使用.来访问数据，代表的是当前作用域的当前对象，类似于 Java 的 this、Python 的 self.
	// 所以.Name 代表的是访问传入的 User 的 Name 字段
	parse, err := tpl.Parse(`Hello, {{index . 0 }}`)
	require.NoError(t, err)
	buffer := &bytes.Buffer{}
	// 指针结构体都行
	err = parse.Execute(buffer, []string{"zhangsan"})
	require.NoError(t, err)
	assert.Equal(t, `Hello, zhangsan`, buffer.String())
}

func TestBase(t *testing.T) {
	tpl := template.New("hello-world")
	// 使用.来访问数据，代表的是当前作用域的当前对象，类似于 Java 的 this、Python 的 self.
	// 所以.Name 代表的是访问传入的 User 的 Name 字段
	parse, err := tpl.Parse(`Hello, {{ . }}`)
	require.NoError(t, err)
	buffer := &bytes.Buffer{}
	// 指针结构体都行
	err = parse.Execute(buffer, 123)
	require.NoError(t, err)
	assert.Equal(t, `Hello, 123`, buffer.String())
}

// 模板基本语法 -- 变量声明
// 去除空格和换行:注意要和别的元素用空格分开声明变量:如同 Go语言，但是用 $ 来表示。$xxx :=some value
// 执行方法调用: 形式“调用者.方法 参数1 参数2”，注意，方法调用的形式和 Go 语言本身的调用形式差异很大
//const serviceTpl = `
//{{- $service :=.GenName -}}
//type {{ $service }} struct {
//	Endpoint string
//	Path string
//	Client http.Client
//
//}
//`

// 方法调用
func TestFuncCall(t *testing.T) {
	tpl := template.New("hello-world")

	parse, err := tpl.Parse(`
切片长度: {{ len .Slice }}
{{printf "%.2f" 1.234}}
Hello, {{ .Hello "Tom" "Jerry"}}`)
	require.NoError(t, err)
	buffer := &bytes.Buffer{}
	// 指针结构体都行
	err = parse.Execute(buffer, FuncCall{
		Slice: []string{"a", "b"},
	})
	require.NoError(t, err)
	assert.Equal(t, `
切片长度: 2
1.23
Hello, Tom·Jerry`, buffer.String())

}

// 模板基本语法- 循环
// 使用 range 关键字，形式 range $idx, $elem := 某个切片在右图中，我们迭代了.Methods注意:不支持 for...i... 的循环形式，也不支持 for true 这种形式
func TestForLoop(t *testing.T) {
	tpl := template.New("hello-world")
	//	parse, err := tpl.Parse(`
	//{{- range $idx, $ele := .}}
	//{{- .}}
	//{{$idx}}-{{$ele}}
	//{{end}}
	//`)
	parse, err := tpl.Parse(`
{{- range $idx, $ele := .}}
{{- $idx}},
{{- end}}`)
	require.NoError(t, err)
	buffer := &bytes.Buffer{}
	// 想要迭代一百次 里面可以任意类型  直接借用下下表而已不关心什么元素
	data := make([]bool, 100)
	//err = parse.Execute(buffer, FuncCall{
	//	Slice: []string{"a", "b"},
	//})
	err = parse.Execute(buffer, data)
	require.NoError(t, err)
	//	assert.Equal(t, `a
	//0-a
	//b
	//1-b
	//
	//`, buffer.String())
	assert.Equal(t, `0,1,2,3,4,5,6,7,8,9,10,11,12,13,14,15,16,17,18,19,20,21,22,23,24,25,26,27,28,29,30,31,32,33,34,35,36,37,38,39,40,41,42,43,44,45,46,47,48,49,50,51,52,53,54,55,56,57,58,59,60,61,62,63,64,65,66,67,68,69,70,71,72,73,74,75,76,77,78,79,80,81,82,83,84,85,86,87,88,89,90,91,92,93,94,95,96,97,98,99,`, buffer.String())

}

// 条件判断
// 一样采用 if-else 或者 if-else if 的结构
// 可以使用 and: and 条件1 条件2
// 可以使用 or: or 条件1 条件2
// 可以使用 not: not 条件1
func TestIfElse(t *testing.T) {

	type User struct {
		Age int
	}

	tpl := template.New("hello-world")
	parse, err := tpl.Parse(`
{{- if and (gt .Age 0) (le .Age 6) }}
儿童: (0, 6]
{{ else if and (gt .Age 6) (le .Age 18) }}
少年: (6, 18]
{{ else }}
成人: >18
{{end -}}
`)
	require.NoError(t, err)
	buffer := &bytes.Buffer{}

	err = parse.Execute(buffer, User{Age: 20})
	require.NoError(t, err)
	assert.Equal(t, `
成人: >18
`, buffer.String())

}

type FuncCall struct {
	Slice []string
}

func (f FuncCall) Hello(first, last string) string {
	return fmt.Sprintf("%s·%s", first, last)
}
