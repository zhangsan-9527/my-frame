package model

import (
	"database/sql"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"my-frame/orm/internal/errs"
	"reflect"
	"sync"
	"testing"
)

func Test_parse_Registry(t *testing.T) {

	tests := []struct {
		name      string
		entity    any
		wantModel *Model
		fields    []*Field
		wantErr   error
	}{
		{
			name:    "test Model",
			entity:  TestModel{},
			wantErr: errs.ErrPointerOnly,
		},
		{
			name:   "pointer",
			entity: &TestModel{},
			wantModel: &Model{
				TableName: "test_model",
			},
			fields: []*Field{
				{
					ColName: "id",
					GoName:  "Id",
					Typ:     reflect.TypeOf(int64(0)),
				},
				{
					ColName: "first_name",
					GoName:  "FirstName",
					Typ:     reflect.TypeOf(""),
				},
				{
					ColName: "last_name",
					GoName:  "LastName",
					Typ:     reflect.TypeOf(&sql.NullString{}),
				},
				{
					ColName: "age",
					GoName:  "Age",
					Typ:     reflect.TypeOf(int8(0)),
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
			m, err := r.Register(tt.entity)
			assert.Equal(t, tt.wantErr, err)
			if err != nil {
				return
			}
			fieldMap := make(map[string]*Field)
			columnMap := make(map[string]*Field)
			for _, f := range tt.fields {
				fieldMap[f.GoName] = f
				columnMap[f.ColName] = f
			}
			tt.wantModel.FieldMap = fieldMap
			tt.wantModel.ColumnMap = columnMap
			assert.Equal(t, tt.wantModel, m)

		})
	}
}

func TestRegistry_get(t *testing.T) {
	tests := []struct {
		name string

		entity    any
		wantModel *Model
		fields    []*Field
		wantErr   error

		cacheSize int
	}{
		{
			name:   "pointer",
			entity: &TestModel{},
			wantModel: &Model{
				TableName: "test_model",
			},
			fields: []*Field{
				{
					ColName: "id",
					GoName:  "Id",
					Typ:     reflect.TypeOf(int64(0)),
				},
				{
					ColName: "first_name",
					GoName:  "FirstName",
					Typ:     reflect.TypeOf(""),
				},
				{
					ColName: "last_name",
					GoName:  "LastName",
					Typ:     reflect.TypeOf(&sql.NullString{}),
				},
				{
					ColName: "age",
					GoName:  "Age",
					Typ:     reflect.TypeOf(int8(0)),
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
			wantModel: &Model{
				TableName: "tag_table",
			},
			fields: []*Field{
				{
					ColName: "first_name_t",
					GoName:  "FirstName",
					Typ:     reflect.TypeOf(""),
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
			wantModel: &Model{
				TableName: "tag_table",
			},
			fields: []*Field{
				{
					ColName: "first_name",
					GoName:  "FirstName",
					Typ:     reflect.TypeOf(""),
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
			wantModel: &Model{
				TableName: "tag_table",
			},
			fields: []*Field{
				{
					ColName: "first_name",
					GoName:  "FirstName",
					Typ:     reflect.TypeOf(""),
				},
			},

			cacheSize: 1,
		},
		{
			name:   "table name",
			entity: &CustomTableName{},
			wantModel: &Model{
				TableName: "custom_table_name_t",
			},
			fields: []*Field{
				{
					ColName: "first_name",
					GoName:  "FirstName",
					Typ:     reflect.TypeOf(""),
				},
			},
			cacheSize: 1,
		},
		{
			name:   "table name",
			entity: &CustomTableNamePtr{},
			wantModel: &Model{
				TableName: "custom_table_name_ptr_t",
			},
			fields: []*Field{
				{
					ColName: "first_name",
					GoName:  "FirstName",
					Typ:     reflect.TypeOf(""),
				},
			},
			cacheSize: 1,
		},
		{
			name:   "empty table name",
			entity: &EmptyTableName{},
			wantModel: &Model{
				TableName: "empty_table_name",
			},
			fields: []*Field{
				{
					ColName: "first_name",
					GoName:  "FirstName",
					Typ:     reflect.TypeOf(""),
				},
			},
			cacheSize: 1,
		},
	}

	r := NewRegistry()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m, err := r.Get(tt.entity)
			if err != nil {
				assert.Equal(t, tt.wantErr, err)
				return
			}

			fieldMap := make(map[string]*Field)
			columnMap := make(map[string]*Field)
			for _, f := range tt.fields {
				fieldMap[f.GoName] = f
				columnMap[f.ColName] = f
			}
			tt.wantModel.FieldMap = fieldMap
			tt.wantModel.ColumnMap = columnMap
			assert.Equal(t, tt.wantModel, m)
			// 只是检测数量
			//assert.Equal(t, tt.cacheSize, getSyncMapLength(&r.models))

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

func TestModelWithTableName(t *testing.T) {
	r := NewRegistry()
	m, err := r.Register(&TestModel{}, ModelWithTableName("test_model_tttt"))
	require.NoError(t, err)
	assert.Equal(t, "test_model_tttt", m.TableName)
}

func TestModelWithCloumnName(t *testing.T) {
	testCases := []struct {
		name    string
		field   string
		colName string

		wantCloName string
		wantErr     error
	}{
		{
			name:        "cloumn name",
			field:       "FirstName",
			colName:     "first_name_cccc",
			wantCloName: "first_name_cccc",
		},
		{
			name:        "invalid cloumn name",
			field:       "XXX",
			colName:     "first_name_cccc",
			wantCloName: "first_name_cccc",
			wantErr:     errs.NewErrUnknownField("XXX"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			r := NewRegistry()
			m, err := r.Register(&TestModel{}, ModelWithColumnName(tc.field, tc.colName))
			assert.Equal(t, tc.wantErr, err)
			if err != nil {
				return
			}
			fd, ok := m.FieldMap[tc.field]
			assert.True(t, ok)
			assert.Equal(t, tc.wantCloName, fd.ColName)
		})
	}
}

type TestModel struct {
	Id        int64
	FirstName string
	Age       int8
	LastName  *sql.NullString
}
