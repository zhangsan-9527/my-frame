package valuer

import (
	"database/sql"
)

type unsafeValue struct {
}

func (u unsafeValue) SetColumns(rows *sql.Rows) error {
	//TODO implement me
	panic("implement me")
}
