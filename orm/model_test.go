package orm

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
				tableName: "test_model",
			},
			fields: []*Field{
				{
					colName: "id",
					goName:  "Id",
					typ:     reflect.TypeOf(int64(0)),
				},
				{
					colName: "first_name",
					goName:  "FirstName",
					typ:     reflect.TypeOf(""),
				},
				{
					colName: "last_name",
					goName:  "LastName",
					typ:     reflect.TypeOf(&sql.NullString{}),
				},
				{
					colName: "age",
					goName:  "Age",
					typ:     reflect.TypeOf(int8(0)),
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
				fieldMap[f.goName] = f
				columnMap[f.colName] = f
			}
			tt.wantModel.fieldMap = fieldMap
			tt.wantModel.columnMap = columnMap
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
				tableName: "test_model",
			},
			fields: []*Field{
				{
					colName: "id",
					goName:  "Id",
					typ:     reflect.TypeOf(int64(0)),
				},
				{
					colName: "first_name",
					goName:  "FirstName",
					typ:     reflect.TypeOf(""),
				},
				{
					colName: "last_name",
					goName:  "LastName",
					typ:     reflect.TypeOf(&sql.NullString{}),
				},
				{
					colName: "age",
					goName:  "Age",
					typ:     reflect.TypeOf(int8(0)),
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
				tableName: "tag_table",
			},
			fields: []*Field{
				{
					colName: "first_name_t",
					goName:  "FirstName",
					typ:     reflect.TypeOf(""),
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
				tableName: "tag_table",
			},
			fields: []*Field{
				{
					colName: "first_name",
					goName:  "FirstName",
					typ:     reflect.TypeOf(""),
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
				tableName: "tag_table",
			},
			fields: []*Field{
				{
					colName: "first_name",
					goName:  "FirstName",
					typ:     reflect.TypeOf(""),
				},
			},

			cacheSize: 1,
		},
		{
			name:   "table name",
			entity: &CustomTableName{},
			wantModel: &Model{
				tableName: "custom_table_name_t",
			},
			fields: []*Field{
				{
					colName: "first_name",
					goName:  "FirstName",
					typ:     reflect.TypeOf(""),
				},
			},
			cacheSize: 1,
		},
		{
			name:   "table name",
			entity: &CustomTableNamePtr{},
			wantModel: &Model{
				tableName: "custom_table_name_ptr_t",
			},
			fields: []*Field{
				{
					colName: "first_name",
					goName:  "FirstName",
					typ:     reflect.TypeOf(""),
				},
			},
			cacheSize: 1,
		},
		{
			name:   "empty table name",
			entity: &EmptyTableName{},
			wantModel: &Model{
				tableName: "empty_table_name",
			},
			fields: []*Field{
				{
					colName: "first_name",
					goName:  "FirstName",
					typ:     reflect.TypeOf(""),
				},
			},
			cacheSize: 1,
		},
	}

	r := newRegistry()
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
				fieldMap[f.goName] = f
				columnMap[f.colName] = f
			}
			tt.wantModel.fieldMap = fieldMap
			tt.wantModel.columnMap = columnMap
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
	r := newRegistry()
	m, err := r.Register(&TestModel{}, ModelWithTableName("test_model_tttt"))
	require.NoError(t, err)
	assert.Equal(t, "test_model_tttt", m.tableName)
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
			r := newRegistry()
			m, err := r.Register(&TestModel{}, ModelWithColumnName(tc.field, tc.colName))
			assert.Equal(t, tc.wantErr, err)
			if err != nil {
				return
			}
			fd, ok := m.fieldMap[tc.field]
			assert.True(t, ok)
			assert.Equal(t, tc.wantCloName, fd.colName)
		})
	}
}
