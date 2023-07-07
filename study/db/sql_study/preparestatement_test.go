package sql_study

import (
	"context"
	"database/sql"
	"github.com/stretchr/testify/require"
	"log"
	"testing"
	"time"
)

// 预编译
func TestPrepareStatement(t *testing.T) {
	//region 初始化db和建表
	db, err := sql.Open("sqlite3", "file:test.db?cache=shared&mode=memory")
	require.NoError(t, err)
	defer db.Close()
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()
	_, err = db.ExecContext(ctx, `
CREATE TABLE IF NOT EXISTS test_model(
    id INTEGER PRIMARY KEY,
    first_name TEXT NOT NULL,
    age INTEGER,
    last_name TEXT NOT NULL
)
`)
	require.NoError(t, err)
	//endregion

	// 预编译返回的statement 可以重复调用, 只需要传递参数值即可
	stmt, err := db.PrepareContext(ctx, "SELECT * FROM `test_model` WHERE `id`=?")
	defer stmt.Close() // stmt不需要的时候需要调用close, 通常在整个应用关闭的时候关闭
	require.NoError(t, err)

	rows, err := stmt.QueryContext(ctx, 1) // id=1
	require.NoError(t, err)
	for rows.Next() {
		tm := TestModel{}
		err = rows.Scan(&tm.Id, &tm.FirstName, &tm.Age, &tm.LastName)
		require.NoError(t, err)
		log.Println(tm)
	}
	stmt.Close()

}
