package file_demo

import (
	"fmt"
	"github.com/stretchr/testify/require"
	"os"
	"testing"
)

func TestFile(t *testing.T) {

	fmt.Println(os.Getwd())

	file, err := os.Open("testdata/my_file.txt") // Open调用OpenFile(name, O_RDONLY, 0)  O_RDONLY表示只读的
	require.NoError(t, err)
	data := make([]byte, 64)
	n, err := file.Read(data)
	fmt.Println(n)
	require.NoError(t, err)

	// 一定会报错
	n, err = file.WriteString("hello")
	fmt.Println(n)
	// bad file description 不可写
	fmt.Println(err)
	require.Error(t, err)
	file.Close()

	f, err := os.OpenFile("testdata/my_file.txt", os.O_APPEND|os.O_WRONLY, os.ModeAppend)
	require.NoError(t, err)
	n, err = f.WriteString("hello")
	fmt.Println(n)
	require.NoError(t, err)
	file.Close()

	f, err = os.Create("testdata/my_file_copy.txt")
	require.NoError(t, err)
	n, err = f.WriteString("hello, world")
	fmt.Println(n)
	require.NoError(t, err)
}
