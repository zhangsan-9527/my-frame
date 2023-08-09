package orm

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"my-frame/orm/internal/errs"
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

func memoryDB(t *testing.T) *DB {
	db, err := sql.Open("sqlite3", "C:\\Users\\81933\\test.db?cache=shared&mode=memory")
	require.NoError(t, err)
	return &DB{
		r:  newRegistry(),
		db: db,
	}
}
