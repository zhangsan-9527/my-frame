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
	db, err := sql.Open("sqlite3", "C:\\Users\\81933\\test.db?cache=shared&mode=memory")
	require.NoError(t, err)
	defer db.Close()
	// 这里 你就可以用 DB 了
	//sql.OpenDB()

	//db.Exec() 与 db.ExecContext() 区别有无 ctx  用于控制超时等
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	//// 除了 SELECT 语句, 都是使用 ExecContext
	//_, err = db.ExecContext(ctx, `
	//CREATE TABLE IF NOT EXISTS test_model(
	//   id INTEGER PRIMARY KEY,
	//   first_name TEXT NOT NULL,
	//   age INTEGER,
	//   last_name TEXT NOT NULL
	//)
	//`)

	// 完成了建表
	require.NoError(t, err)

	// 使用 ? 作为查询的参数的占位符 防止依赖注入
	res, err := db.ExecContext(ctx, "INSERT INTO test_model (`id`, `first_name`, `age`, `last_name` ) VALUES (?,?,?,?)", 1, "zs", 18, "9527")

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

	// 预期至少会有一行所以如果没有值调用Scan 会有sql.ErrNoRows 但是 QueryContext 没有
	row = db.QueryRowContext(ctx, "SELECT `id`, `first_name`, `age`, `last_name` FROM `test_model` WHERE `id` = ?", 2)
	// 查询不到
	require.NoError(t, row.Err())
	err = row.Scan(&tm.Id, &tm.FirstName, &tm.Age, &tm.LastName)
	// 在这里返回 err
	require.Error(t, sql.ErrNoRows, err)

	// 需要调用Next() 批量查询
	rows, err := db.QueryContext(ctx, "SELECT `id`, `first_name`, `age`, `last_name` FROM `test_model` WHERE `id` = ?", 2)

	for rows.Next() {
		tm = TestModel{}
		// 这里不会返回sql.ErrNoRows
		err = rows.Scan(&tm.Id, &tm.FirstName, &tm.Age, &tm.LastName)
		log.Println(tm)
		require.NoError(t, err)
	}

	cancel()
}

func TestTx(t *testing.T) {
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
	require.NoError(t, err)

	/*
		事务 API-_TxOptions

			大多数时候我们不需要设置这个
			两个参数: ReadOnly 只读 , lsolation 隔离级别
			逼不得已要设置lsolation 字段的时候，要确认自己使用的数据库支持该级别，并且弄清楚效果
			(同一个隔离级别在不同的数据库上都有不同的解释要小心使用)
	*/

	tx, err := db.BeginTx(ctx, &sql.TxOptions{})
	require.NoError(t, err)

	// 使用 ? 作为查询的参数的占位符 防止依赖注入
	res, err := tx.ExecContext(ctx, "INSERT INTO test_model (`id`, `first_name`, `age`, `last_name` ) VALUES (?,?,?,?), 1, 'zs', 18, '9527'")
	if err != nil {
		// 回滚事务
		err = tx.Rollback()
		if err != nil {
			log.Println(err)
		}
	}

	require.NoError(t, err)
	affected, err := res.RowsAffected() // 受影响影响行数
	require.NoError(t, err)
	log.Println("影响行数:", affected)

	id, err := res.LastInsertId() // 最后插入的id
	require.NoError(t, err)
	log.Println("最后插入 ID:", id)

	// 提交事务
	err = tx.Commit()

	cancel()

}

type TestModel struct {
	Id        int64
	FirstName string
	Age       int8
	LastName  *sql.NullString
}
