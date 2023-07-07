package sql_study

import (
	"context"
	"errors"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/require"
	"log"
	"testing"
)

/*
在单元测试里面我们不希望依赖于真实的数据库，
因为数据难以模拟， 而且 error 更加难以模拟，
所以我们采用 sqlmock 来做单元测试。
sqlmock 使用：
• 初始化： 返回一个 mockDB， 类型是*sql.DB， 还有 mock 用于构造模拟的场景；
• 设置 mock： 基本上是 ExpectXXX  WillXXX， 严格依赖于顺序
*/

func TestSQLMock(t *testing.T) {

	db, mock, err := sqlmock.New()
	defer db.Close()
	require.NoError(t, err)

	// 构造模拟数据
	// 注意mock.ExpectXXX跟 db.Query/Exec 必须是一一对应的.
	// 第一个
	mockRows := sqlmock.NewRows([]string{"id", "first_name"})
	mockRows.AddRow(1, "Tom")
	// 正则表达式
	mock.ExpectQuery("SELECT id,first_name FROM `user`.*").WillReturnRows(mockRows)       //第一个ExpectQuery
	mock.ExpectQuery("SELECT id FROM `user`.*").WillReturnError(errors.New("mock error")) //第二个ExpectQuery

	// result :=sqlmock.NewResult()
	// mock.ExpectExec().WillReturnResult()

	rows, err := db.QueryContext(context.Background(), "SELECT id,first_name FROM `user` WHERE id=1") //第一个QueryContext
	require.NoError(t, err)
	for rows.Next() {
		tm := TestModel{}
		err = rows.Scan(&tm.Id, &tm.FirstName)
		require.NoError(t, err)
		log.Println(tm)
	}

	_, err = db.QueryContext(context.Background(), "SELECT id FROM `user` WHERE id=1") //第二个QueryContext
	require.Error(t, err)
}
