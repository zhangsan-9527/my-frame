package reflect

import (
	"errors"
	"github.com/stretchr/testify/assert"
	"reflect"
	"testing"
)

func TestIterateField(t *testing.T) {

	type User struct {
		Name string
		age  int
	}

	testCases := []struct {
		name string

		entity any

		wantErr error
		wantRes map[string]any
	}{
		{
			name: "struct",
			entity: User{
				Name: "Tom",
				age:  18,
			},
			wantRes: map[string]any{
				"Name": "Tom",
				// age 私有的, 拿不到,最终我们创建了0来填充
				"age": 0,
			},
		},
		{
			name: "pointer",
			entity: &User{
				Name: "Tom",
				age:  18,
			},
			wantRes: map[string]any{
				"Name": "Tom",
				// age 私有的, 拿不到,最终我们创建了0来填充
				"age": 0,
			},
		},
		{
			name: "multiple pointer",
			entity: func() **User { // 二级指针
				res := &User{
					Name: "Tom",
					age:  18,
				}
				return &res
			}(),
			wantRes: map[string]any{
				"Name": "Tom",
				// age 私有的, 拿不到,最终我们创建了0来填充
				"age": 0,
			},
		},
		{
			name:    "base type",
			entity:  18,
			wantErr: errors.New("不支持类型"),
		},
		{
			name:    "nil",
			entity:  nil,
			wantErr: errors.New("不支持 nil"),
		},
		{
			name:    "user nil",
			entity:  (*User)(nil),
			wantErr: errors.New("不支持零值"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			res, err := IterateField(tc.entity)
			assert.Equal(t, tc.wantErr, err)
			if err != nil {
				return
			}
			assert.Equal(t, tc.wantRes, res)
		})
	}
}

func TestSetField(t *testing.T) {
	type User struct {
		Name string
		age  int
	}

	testCases := []struct {
		name string

		entity   any
		field    string
		newValue any

		wantErr error
		// 修改后的 entity
		wantEntity any
	}{
		{
			name: "struct",
			entity: User{
				Name: "zhangsan",
			},
			field:    "Name",
			newValue: "lisi",
			wantErr:  errors.New("不可修改字段"),
		},
		{
			name: "pointer",
			entity: &User{
				Name: "zhangsan",
			},
			field:      "Name",
			newValue:   "lisi",
			wantEntity: &User{Name: "lisi"},
		},
		{
			name: "pointer exported",
			entity: &User{
				age: 18,
			},
			field:    "age",
			newValue: 19,
			wantErr:  errors.New("不可修改字段"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := SetField(tc.entity, tc.field, tc.newValue)
			assert.Equal(t, tc.wantErr, err)
			if err != nil {
				return
			}
			assert.Equal(t, tc.wantEntity, tc.entity)
		})
	}

	var i = 0
	ptr := &i
	reflect.ValueOf(ptr).Elem().Set(reflect.ValueOf(12))
	assert.Equal(t, i, 12)
}
