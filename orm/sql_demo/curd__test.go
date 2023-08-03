package sql_demo

import (
	"context"
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
	"github.com/stretchr/testify/require"
	"log"
	"testing"
	"time"
)

func TestDB(t *testing.T) {
	db, err := sql.Open("sqlite3", "file:test.db?cache=shared&mode=memory")
	require.NoError(t, err)
	defer db.Close()
	// 这里 你就可以用 DB 了
	//sql.OpenDB()

	//db.Exec() 与 db.ExecContext() 区别有无 ctx  用于控制超时等
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	// 除了 SELECT 语句, 都是使用 ExecContext
	_, err = db.ExecContext(ctx, `
CREATE TABLE IF NOT EXISTS test_model(
    id INTEGER PRIMARY KEY,
    first_name TEXT NOT NULL,
    age INTEGER,
    last_name TEXT NOT NULL
)
`)

	// 完成了建表
	require.NoError(t, err)

	// 使用 ? 作为查询的参数的占位符 防止依赖注入
	res, err := db.ExecContext(ctx, "INSERT INTO test_model (`id`, `first_name`, `age`, `last_name` ) VALUES (?,?,?,?), 1, 'zs', 18, '9527'")

	require.NoError(t, err)
	affected, err := res.RowsAffected() // 受影响影响行数
	require.NoError(t, err)
	log.Println(affected)

	id, err := res.LastInsertId() // 最后插入的id
	require.NoError(t, err)
	log.Println(id)

	// 返回一行
	row := db.QueryRowContext(ctx, "SELECT `id`, `first_name`, `age`, `last_name` FROM `test_model` WHERE `id` = ?", 1)
	require.NoError(t, row.Err())
	tm := TestModel{}
	err = row.Scan(&tm.Id, &tm.FirstName, &tm.Age, &tm.LastName)
	require.NoError(t, err)

	row = db.QueryRowContext(ctx, "SELECT `id`, `first_name`, `age`, `last_name` FROM `test_model` WHERE `id` = ?", 2)
	// 查询不到
	require.NoError(t, row.Err())
	err = row.Scan(&tm.Id, &tm.FirstName, &tm.Age, &tm.LastName)
	// 在这里返回 err
	require.Error(t, sql.ErrNoRows, err)

	cancel()
}

type TestModel struct {
	Id        int64
	FirstName string
	Age       int8
	LastName  *sql.NullString
}
