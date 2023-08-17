package orm

import (
	"context"
	"database/sql"
	"errors"
	"github.com/DATA-DOG/go-sqlmock"
	_ "github.com/mattn/go-sqlite3"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"my-frame/orm/internal/errs"
	"my-frame/orm/model"
	"testing"
)

func TestSelector_Build(t *testing.T) {
	db := memoryDB(t)
	testCass := []struct {
		name      string
		bulider   QueryBuilder
		wantQuery *Query
		wantErr   error
	}{
		{
			name:    "no from",
			bulider: NewSelector[TestModel](db),
			wantQuery: &Query{
				SQL:  "SELECT * FROM `test_model`;",
				Args: nil,
			},
		},
		{
			name:    "from",
			bulider: NewSelector[TestModel](db).From("test_model"),
			wantQuery: &Query{
				SQL:  "SELECT * FROM test_model;",
				Args: nil,
			},
		},
		{
			name:    "empty from",
			bulider: NewSelector[TestModel](db).From(""),
			wantQuery: &Query{
				SQL:  "SELECT * FROM `test_model`;",
				Args: nil,
			},
		},
		{
			name:    "with db",
			bulider: NewSelector[TestModel](db).From("test_db.test_model"),
			wantQuery: &Query{
				SQL:  "SELECT * FROM test_db.test_model;",
				Args: nil,
			},
		},
		{
			name:    "where",
			bulider: NewSelector[TestModel](db).Where(C("Age").Eq(18)),
			wantQuery: &Query{
				SQL:  "SELECT * FROM `test_model` WHERE `age` = ?;",
				Args: []any{18},
			},
		},
		{
			name:    "not",
			bulider: NewSelector[TestModel](db).Where(Not(C("Age").Eq(18))),
			wantQuery: &Query{
				SQL:  "SELECT * FROM `test_model` WHERE  NOT (`age` = ?);",
				Args: []any{18},
			},
		},
		{
			name:    "and",
			bulider: NewSelector[TestModel](db).Where(C("Age").Eq(18).And(C("Age").Eq(30))),
			wantQuery: &Query{
				SQL:  "SELECT * FROM `test_model` WHERE (`age` = ?) AND (`age` = ?);",
				Args: []any{18, 30},
			},
		},
		{
			name:    "or",
			bulider: NewSelector[TestModel](db).Where(C("Age").Eq(18).Or(C("Age").Eq(30))),
			wantQuery: &Query{
				SQL:  "SELECT * FROM `test_model` WHERE (`age` = ?) OR (`age` = ?);",
				Args: []any{18, 30},
			},
		},
		{
			name:    "invalid column",
			bulider: NewSelector[TestModel](db).Where(C("Age").Eq(18).Or(C("XXXX").Eq(30))),
			wantErr: errs.NewErrUnknownField("XXXX"),
		},
	}

	for _, tc := range testCass {
		t.Run(tc.name, func(t *testing.T) {
			query, err := tc.bulider.Build()
			assert.Equal(t, tc.wantErr, err)
			if err != nil {
				return
			}
			assert.Equal(t, tc.wantQuery, query)
		})
	}

}

type TestModel struct {
	Id        int64
	FirstName string
	Age       int8
	LastName  *sql.NullString
}

func TestSelector_Get(t *testing.T) {
	mockDB, mock, err := sqlmock.New()
	require.NoError(t, err)
	db, err := OpenDB(mockDB)
	require.NoError(t, err)

	// 对应于 query error
	mock.ExpectQuery("SELECT .*").WillReturnError(errors.New("query error"))

	// 对应于 no rows
	rows := sqlmock.NewRows([]string{"id", "first_name", "age", "last_name"})
	mock.ExpectQuery("SELECT .*").WillReturnRows(rows)

	//// scan error
	//rows = sqlmock.NewRows([]string{"id", "first_name", "age", "last_name"})
	//rows.AddRow("ABC", "Tom", "18", "Henry")
	//mock.ExpectQuery("SELECT .*").WillReturnRows(rows)

	// 对应于 data
	rows = sqlmock.NewRows([]string{"id", "first_name", "age", "last_name"})
	rows.AddRow("1", "Tom", "18", "Henry")
	mock.ExpectQuery("SELECT .*").WillReturnRows(rows)

	testCases := []struct {
		name    string
		s       *Selector[TestModel]
		wantErr error
		wantRes *TestModel
	}{
		{
			name:    "invalid query",
			s:       NewSelector[TestModel](db).Where(C("XXX").Eq(1)),
			wantErr: errs.NewErrUnknownField("XXX"),
		},
		{
			name:    "query error",
			s:       NewSelector[TestModel](db).Where(C("Id").Eq(1)),
			wantErr: errors.New("query error"),
		},
		{
			name:    "no rows",
			s:       NewSelector[TestModel](db).Where(C("Id").LT(1)),
			wantErr: ErrNoRows,
		},
		//{
		//	name:    "scan error",
		//	s:       NewSelector[TestModel](db).Where(C("Id").LT(1)),
		//	wantErr: ErrNoRows,
		//},
		{
			name: "data",
			s:    NewSelector[TestModel](db).Where(C("Id").LT(1)),
			wantRes: &TestModel{
				Id:        1,
				FirstName: "Tom",
				Age:       18,
				LastName:  &sql.NullString{Valid: true, String: "Henry"},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			res, err := tc.s.Get(context.Background())
			require.Equal(t, tc.wantErr, err)
			if err != nil {
				return
			}
			assert.Equal(t, tc.wantRes, res)
		})
	}
}

func TestSelector_GetV1(t *testing.T) {
	mockDB, mock, err := sqlmock.New()
	require.NoError(t, err)
	db, err := OpenDB(mockDB)
	require.NoError(t, err)

	// 对应于 query error
	mock.ExpectQuery("SELECT .*").WillReturnError(errors.New("query error"))

	// 对应于 no rows
	rows := sqlmock.NewRows([]string{"id", "first_name", "age", "last_name"})
	mock.ExpectQuery("SELECT .*").WillReturnRows(rows)

	//// scan error
	//rows = sqlmock.NewRows([]string{"id", "first_name", "age", "last_name"})
	//rows.AddRow("ABC", "Tom", "18", "Henry")
	//mock.ExpectQuery("SELECT .*").WillReturnRows(rows)

	// 对应于 data
	rows = sqlmock.NewRows([]string{"id", "first_name", "age", "last_name"})
	rows.AddRow("1", "Tom", "18", "Henry")
	mock.ExpectQuery("SELECT .*").WillReturnRows(rows)

	testCases := []struct {
		name    string
		s       *Selector[TestModel]
		wantErr error
		wantRes *TestModel
	}{
		{
			name:    "invalid query",
			s:       NewSelector[TestModel](db).Where(C("XXX").Eq(1)),
			wantErr: errs.NewErrUnknownField("XXX"),
		},
		{
			name:    "query error",
			s:       NewSelector[TestModel](db).Where(C("Id").Eq(1)),
			wantErr: errors.New("query error"),
		},
		{
			name:    "no rows",
			s:       NewSelector[TestModel](db).Where(C("Id").LT(1)),
			wantErr: ErrNoRows,
		},
		//{
		//	name:    "scan error",
		//	s:       NewSelector[TestModel](db).Where(C("Id").LT(1)),
		//	wantErr: ErrNoRows,
		//},
		{
			name: "data",
			s:    NewSelector[TestModel](db).Where(C("Id").LT(1)),
			wantRes: &TestModel{
				Id:        1,
				FirstName: "Tom",
				Age:       18,
				LastName:  &sql.NullString{Valid: true, String: "Henry"},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			res, err := tc.s.Get(context.Background())
			require.Equal(t, tc.wantErr, err)
			if err != nil {
				return
			}
			assert.Equal(t, tc.wantRes, res)
		})
	}
}

func memoryDB(t *testing.T) *DB {
	db, err := sql.Open("sqlite3", "C:\\Users\\81933\\test.db?cache=shared&mode=memory")
	require.NoError(t, err)
	return &DB{
		r:  model.NewRegistry(),
		db: db,
	}
}
