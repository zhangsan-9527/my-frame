package orm

import (
	"github.com/stretchr/testify/assert"
	"my-frame/orm/internal/errs"
	"reflect"
	"sync"
	"testing"
)

func Test_parseModel(t *testing.T) {

	tests := []struct {
		name      string
		entity    any
		wantModel *model
		wantErr   error
	}{
		{
			name:    "test model",
			entity:  TestModel{},
			wantErr: errs.ErrPointerOnly,
		},
		{
			name:   "pointer",
			entity: &TestModel{},
			wantModel: &model{
				tableName: "test_model",
				fields: map[string]*field{
					"Id": {
						colName: "id",
					},
					"FirstName": {
						colName: "first_name",
					},
					"LastName": {
						colName: "last_name",
					},
					"Age": {
						colName: "age",
					},
				},
			},
		},
		{
			name:    "map",
			entity:  map[string]string{},
			wantErr: errs.ErrPointerOnly,
		},
		{
			name:    "slice",
			entity:  []int{},
			wantErr: errs.ErrPointerOnly,
		},
		{
			name:    "base",
			entity:  0,
			wantErr: errs.ErrPointerOnly,
		},
	}
	r := &registry{}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m, err := r.parseModel(tt.entity)
			assert.Equal(t, tt.wantErr, err)
			if err != nil {
				return
			}
			assert.Equal(t, tt.wantModel, m)

		})
	}
}

func TestRegistry_get(t *testing.T) {
	tests := []struct {
		name string

		entity    any
		wantModel *model
		wantErr   error

		cacheSize int
	}{
		{
			name:   "pointer",
			entity: &TestModel{},
			wantModel: &model{
				tableName: "test_model",
				fields: map[string]*field{
					"Id": {
						colName: "id",
					},
					"FirstName": {
						colName: "first_name",
					},
					"LastName": {
						colName: "last_name",
					},
					"Age": {
						colName: "age",
					},
				},
			},
			cacheSize: 1,
		},
		{
			name: "tag",
			entity: func() any {
				type TagTable struct {
					FirstName string `orm:"column=first_name_t"`
				}
				return &TagTable{}
			}(),
			wantModel: &model{
				tableName: "tag_table",
				fields: map[string]*field{
					"FirstName": {
						colName: "first_name_t",
					},
				},
			},
			cacheSize: 1,
		},
		{
			name: "empty column",
			entity: func() any {
				type TagTable struct {
					FirstName string `orm:"column="`
				}
				return &TagTable{}
			}(),
			wantModel: &model{
				tableName: "tag_table",
				fields: map[string]*field{
					"FirstName": {
						colName: "first_name",
					},
				},
			},
			cacheSize: 1,
		},
		{
			name: "column only",
			entity: func() any {
				type TagTable struct {
					FirstName string `orm:"column"`
				}
				return &TagTable{}
			}(),
			wantErr: errs.NewErrInvalidTagContent("column"),
		},
		{
			name: "ignore tag",
			entity: func() any {
				type TagTable struct {
					FirstName string `orm:"abc=abc"`
				}
				return &TagTable{}
			}(),
			wantModel: &model{
				tableName: "tag_table",
				fields: map[string]*field{
					"FirstName": {
						colName: "first_name",
					},
				},
			},
			cacheSize: 1,
		},
		{
			name:   "table name",
			entity: &CustomTableName{},
			wantModel: &model{
				tableName: "custom_table_name_t",
				fields: map[string]*field{
					"FirstName": {
						colName: "first_name",
					},
				},
			},
			cacheSize: 1,
		},
		{
			name:   "table name",
			entity: &CustomTableNamePtr{},
			wantModel: &model{
				tableName: "custom_table_name_ptr_t",
				fields: map[string]*field{
					"FirstName": {
						colName: "first_name",
					},
				},
			},
			cacheSize: 1,
		},
		{
			name:   "empty table name",
			entity: &EmptyTableName{},
			wantModel: &model{
				tableName: "empty_table_name",
				fields: map[string]*field{
					"FirstName": {
						colName: "first_name",
					},
				},
			},
			cacheSize: 1,
		},
	}

	r := newRegistry()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m, err := r.get(tt.entity)
			if err != nil {
				assert.Equal(t, tt.wantErr, err)
				return
			}
			assert.Equal(t, tt.wantModel, m)
			// 只是检测数量
			assert.Equal(t, tt.cacheSize, getSyncMapLength(&r.models))

			typ := reflect.TypeOf(tt.entity)
			cache, ok := r.models.Load(typ)
			assert.True(t, ok)
			assert.Equal(t, tt.wantModel, cache)
		})
	}
}

func getSyncMapLength(m *sync.Map) int {
	length := 0
	m.Range(func(_, _ interface{}) bool {
		length++
		return true
	})
	return length
}

type CustomTableName struct {
	FirstName string
}

func (c CustomTableName) TableName() string {
	return "custom_table_name_t"
}

type CustomTableNamePtr struct {
	FirstName string
}

func (c *CustomTableNamePtr) TableName() string {
	return "custom_table_name_ptr_t"
}

type EmptyTableName struct {
	FirstName string
}

func (c EmptyTableName) TableName() string {
	return ""
}
