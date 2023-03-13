package orm

import (
	"database/sql"
	"geektime-go-study/orm/internal/errs"
	"github.com/stretchr/testify/assert"
	"testing"
)

type TestModel struct {
	Id int64
	// ""
	FirstName string
	Age       int8
	LastName  *sql.NullString
}

func TestSelector_Build(t *testing.T) {
	newTestModel := NewSelector[TestModel]

	testCases := []struct {
		name      string
		q         QueryBuilder
		wantQuery *Query
		wantErr   error
	}{
		{
			name: "no from",
			q:    newTestModel(),
			wantQuery: &Query{
				SQL: "SELECT * FROM `test_model`;",
			},
			wantErr: nil,
		},
		{
			name: "with from",
			q:    newTestModel().From("test_model_t"),
			wantQuery: &Query{
				SQL: "SELECT * FROM test_model_t;",
			},
		},
		{
			// 调用 FROM，但是传入空字符串
			name: "Empty from",
			q:    NewSelector[TestModel]().From(""),
			wantQuery: &Query{
				SQL: "SELECT * FROM `test_model`;",
			},
		},
		{
			// 调用 FROM，同时传入db
			name: "with db From",
			q:    NewSelector[TestModel]().From("`test_db.test_model`"),
			wantQuery: &Query{
				SQL: "SELECT * FROM `test_db.test_model`;",
			},
		},
		{
			// 单一简单条件
			name: "single and simple predicate",
			q: NewSelector[TestModel]().From("`test_model_t`").
				Where(C("Id").EQ(1)),
			wantQuery: &Query{
				SQL:  "SELECT * FROM `test_model_t` WHERE `id` = ?;",
				Args: []any{1},
			},
		},
		{
			// 多个 predicate
			name: "multiple predicates",
			q: NewSelector[TestModel]().
				Where(C("Age").GT(18), C("Age").LT(35)),
			wantQuery: &Query{
				SQL:  "SELECT * FROM `test_model` WHERE (`age` > ?) AND (`age` < ?);",
				Args: []any{18, 35},
			},
		},
		{
			// 使用 AND
			name: "and",
			q: NewSelector[TestModel]().
				Where(C("Age").GT(18).And(C("Age").LT(35))),
			wantQuery: &Query{
				SQL:  "SELECT * FROM `test_model` WHERE (`age` > ?) AND (`age` < ?);",
				Args: []any{18, 35},
			},
		},
		{
			// 使用 OR
			name: "or",
			q: NewSelector[TestModel]().
				Where(C("Age").GT(18).Or(C("Age").LT(35))),
			wantQuery: &Query{
				SQL:  "SELECT * FROM `test_model` WHERE (`age` > ?) OR (`age` < ?);",
				Args: []any{18, 35},
			},
		},
		{
			// 使用 NOT
			name: "not",
			q:    NewSelector[TestModel]().Where(Not(C("Age").GT(18))),
			wantQuery: &Query{
				// NOT 前面有两个空格，因为我们没有对 NOT 进行特殊处理
				SQL:  "SELECT * FROM `test_model` WHERE  NOT (`age` > ?);",
				Args: []any{18},
			},
		},

		{
			name:    "invalid column",
			q:       NewSelector[TestModel]().Where(Not(C("Unkown").GT(18))),
			wantErr: errs.NewErrUnknownField("Unkown"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			res, err := tc.q.Build()
			assert.Equal(t, tc.wantErr, err)
			if err != nil {
				return
			}
			assert.Equal(t, tc.wantQuery, res)

		})
	}

}
