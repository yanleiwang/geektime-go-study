package sql_study

import (
	"context"
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
	"github.com/stretchr/testify/require"
	"log"
	"testing"
	"time"
)

type TestModel struct {
	Id        int64
	FirstName string
	Age       int8
	LastName  *sql.NullString
}

func TestCRUD(t *testing.T) {

	//region 初始化db
	/*
		Open：
		• driver： 也就是驱动的名字， 例如 “mysql”、“sqlite3”
		• dsn： 简单理解就是数据库链接信息
		• 常见错误： 忘记匿名引入 driver 包  _ "github.com/mattn/go-sqlite3"
		OpenDB： 一般用于接入一些自定义的驱动， 例如
		说将分库分表做成一个驱动
		go sql driver list : https://github.com/golang/go/wiki/SQLDrivers

	*/
	db, err := sql.Open("sqlite3", "file:test.db?cache=shared&mode=memory")
	// sql.OpenDB()
	require.NoError(t, err)
	defer db.Close()

	// context用来 超时控制
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	//endregion

	//region 建表操作
	/*
		除了 SELECT 语句，增删改/ddl语句等等都是使用 ExecContext/ Exec
	*/
	{
		_, err = db.ExecContext(ctx, `
CREATE TABLE IF NOT EXISTS test_model(
    id INTEGER PRIMARY KEY,
    first_name TEXT NOT NULL,
    age INTEGER,
    last_name TEXT NOT NULL
)
`)
		require.NoError(t, err)
	}
	//endregion

	//region 增删改
	/*
		增改删：
		• Exec 或 ExecContext
		• 可以用 ExecContext 来控制超时, 所以一般用ExecContext
		• 返回值 sql.Result和error
		sql.Result 支持 RowsAffected() 和 LastInsertId() 分别返回受影响行数和最后插入id
	*/
	{
		// 使用 ？ 作为查询的参数的占位符
		res, err := db.ExecContext(ctx, "INSERT INTO test_model(`id`, `first_name`, `age`, `last_name`) VALUES(?, ?, ?, ?)",
			1, "Tom", 18, "Jerry")
		require.NoError(t, err)
		affected, err := res.RowsAffected()
		require.NoError(t, err)
		log.Println("受影响行数", affected)
		lastId, err := res.LastInsertId()
		require.NoError(t, err)
		log.Println("最后插入 ID", lastId)
	}
	//endregion

	//region 查询
	/*
		查询：
		• QueryRow 和 QueryRowContext： 查询单行数据
		• Query 和 QueryContext： 查询多行数据
		一般用 xxxContext (为了超时控制)
		要注意参数传递， 一般的 SQL 都是使用 ? 作
		为参数占位符。不要把参数拼接进去 SQL 本身，容易引起sql注入
	*/
	{
		// 查询单行数据, 如果没有数据在调用 Row 的 Scan 的时候会返回 sql.ErrNoRow。
		row := db.QueryRowContext(ctx,
			"SELECT `id`, `first_name`, `age`, `last_name` FROM `test_model` WHERE `id` = ?", 1)
		require.NoError(t, row.Err())
		tm := TestModel{}
		err = row.Scan(&tm.Id, &tm.FirstName, &tm.Age, &tm.LastName)
		require.NoError(t, err)

		row = db.QueryRowContext(ctx,
			"SELECT `id`, `first_name`, `age`, `last_name` FROM `test_model` WHERE `id` = ?", 2)
		require.NoError(t, row.Err())
		tm = TestModel{}
		err = row.Scan(&tm.Id, &tm.FirstName, &tm.Age, &tm.LastName)
		require.Error(t, sql.ErrNoRows, err)

		// 查询多行数据,  在调用scan之前 需要调用 .next
		rows, err := db.QueryContext(ctx,
			"SELECT `id`, `first_name`, `age`, `last_name` FROM `test_model` WHERE `id` = ?", 1)
		require.NoError(t, row.Err())
		for rows.Next() {
			tm = TestModel{}
			err = rows.Scan(&tm.Id, &tm.FirstName, &tm.Age, &tm.LastName)
			require.NoError(t, err)
			log.Println(tm)
		}

	}
	//endregion
	cancel()
}
