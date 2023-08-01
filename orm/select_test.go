package orm

import (
	"database/sql"
	"github.com/stretchr/testify/assert"
	"my-frame/orm/internal/errs"
	"testing"
)

func TestSelector_Build(t *testing.T) {

	testCass := []struct {
		name string

		bulider QueryBuilder

		wantQuery *Query
		wantErr   error
	}{
		{
			name:    "no from",
			bulider: &Selector[TestModel]{},
			wantQuery: &Query{
				SQL:  "SELECT * FROM `test_model`;",
				Args: nil,
			},
		},
		{
			name:    "from",
			bulider: (&Selector[TestModel]{}).From("test_model"),
			wantQuery: &Query{
				SQL:  "SELECT * FROM test_model;",
				Args: nil,
			},
		},
		{
			name:    "empty from",
			bulider: (&Selector[TestModel]{}).From(""),
			wantQuery: &Query{
				SQL:  "SELECT * FROM `test_model`;",
				Args: nil,
			},
		},
		{
			name:    "with db",
			bulider: (&Selector[TestModel]{}).From("test_db.test_model"),
			wantQuery: &Query{
				SQL:  "SELECT * FROM test_db.test_model;",
				Args: nil,
			},
		},
		{
			name:    "where",
			bulider: (&Selector[TestModel]{}).Where(C("Age").Eq(18)),
			wantQuery: &Query{
				SQL:  "SELECT * FROM `test_model` WHERE `age` = ?;",
				Args: []any{18},
			},
		},
		{
			name:    "not",
			bulider: (&Selector[TestModel]{}).Where(Not(C("Age").Eq(18))),
			wantQuery: &Query{
				SQL:  "SELECT * FROM `test_model` WHERE  NOT (`age` = ?);",
				Args: []any{18},
			},
		},
		{
			name:    "and",
			bulider: (&Selector[TestModel]{}).Where(C("Age").Eq(18).And(C("Age").Eq(30))),
			wantQuery: &Query{
				SQL:  "SELECT * FROM `test_model` WHERE (`age` = ?) AND (`age` = ?);",
				Args: []any{18, 30},
			},
		},
		{
			name:    "or",
			bulider: (&Selector[TestModel]{}).Where(C("Age").Eq(18).Or(C("Age").Eq(30))),
			wantQuery: &Query{
				SQL:  "SELECT * FROM `test_model` WHERE (`age` = ?) OR (`age` = ?);",
				Args: []any{18, 30},
			},
		},
		{
			name:    "invalid column",
			bulider: (&Selector[TestModel]{}).Where(C("Age").Eq(18).Or(C("XXXX").Eq(30))),
			wantErr: errs.NewErrUnkonwField("XXXX"),
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
