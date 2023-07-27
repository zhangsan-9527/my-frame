package reflect

import (
	"github.com/stretchr/testify/assert"
	"my-frame/orm/reflect/types"
	"reflect"
	"testing"
)

// 如果定义的为结构体 只能访问结构体的方法
// 如果定义的为指针, 不仅能访问指针上的方法还能访问结构体上的方法  输入的第一个参数也是指针
func TestIterateFunc(t *testing.T) {
	type args struct {
		entity any
	}
	tests := []struct {
		name   string
		entity any

		wantRes map[string]FuncInfo
		wantErr error
	}{
		{
			name:   "struct",
			entity: types.NewUser("zhangsan", 18),
			wantRes: map[string]FuncInfo{
				"GetAge": {
					Name: "GetAge",

					// 为什么是types.User{} 看types user.go
					// 下标 0 的指向接收器
					InputTypes:  []reflect.Type{reflect.TypeOf(types.User{})},
					OutputTypes: []reflect.Type{reflect.TypeOf(0)},
					Result:      []any{18},
				},
				//"ChangeName": {
				//	Name:       "ChangeName",
				//	InputTypes: []reflect.Type{reflect.TypeOf("")},
				//	//OutputTypes: []reflect.Type{},
				//	//Result:      []any{},
				//},
			},
		},
		{
			name:   "pointer",
			entity: types.NewUserPtr("zhangsan", 18),
			wantRes: map[string]FuncInfo{
				"GetAge": {
					Name: "GetAge",

					// 为什么是types.User{} 看types user.go
					// 下标 0 的指向接收器
					InputTypes:  []reflect.Type{reflect.TypeOf(&types.User{})},
					OutputTypes: []reflect.Type{reflect.TypeOf(0)},
					Result:      []any{18},
				},
				"ChangeName": {
					Name:        "ChangeName",
					InputTypes:  []reflect.Type{reflect.TypeOf(&types.User{}), reflect.TypeOf("")},
					OutputTypes: []reflect.Type{},
					Result:      []any{},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			res, err := IterateFunc(tt.entity)
			assert.Equal(t, tt.wantErr, err)
			if err != nil {
				return
			}
			assert.Equal(t, tt.wantRes, res)
		})
	}
}

/*
	方法只读没有暴露重写接口

	注意事项:
		方法接收器
			以结构体作为输入，那么只能访问到结构体作为接收器的方法
			以指针作为输入，那么能访问到任何接收器的方法
		输入的第一个参数，永远都是接收器本身

*/
