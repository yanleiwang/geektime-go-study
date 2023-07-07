package sql_study

import (
	"context"
	"database/sql"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"github.com/stretchr/testify/require"
	"log"
	"testing"
	"time"
)

// 事务
func TestTx(t *testing.T) {

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

	/*
		事务API
		• Begin 和 BeginTx
		• Commit
		• Rollback
	*/

	//开启事务,  opt可以配置 是否是可读事务/ 隔离级别
	// 大多数时候我们不需要设置这个
	tx, err := db.BeginTx(ctx, nil)
	require.NoError(t, err)

	res, err := tx.ExecContext(ctx, "INSERT INTO test_model(`id`, `first_name`, `age`, `last_name`) VALUES(?, ?, ?, ?)",
		1, "Tom", 18, "Jerry")
	if err != nil {
		// 回滚事务
		err = tx.Rollback()
		if err != nil {
			log.Println(err)
		}
		return
	}
	require.NoError(t, err)
	fmt.Printf("res: %v\n", res)
	// 提交事务
	err = tx.Commit()

}
