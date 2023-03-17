package orm

import (
	"database/sql"
	"geektime-go-study/orm/internal/errs"
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_registry_get(t *testing.T) {
	type TestModel struct {
		Id int64
		// ""
		FirstName string
		Age       int8
		LastName  *sql.NullString
	}

	testCases := []struct {
		name      string
		val       any
		wantModel *model
		wantErr   error
	}{
		{
			// 指针
			name: "pointer",
			val:  &TestModel{},
			wantModel: &model{
				tableName: "test_model",
				fieldMap: map[string]*field{
					"Id": {
						colName: "id",
					},
					"FirstName": {
						colName: "first_name",
					},
					"Age": {
						colName: "age",
					},
					"LastName": {
						colName: "last_name",
					},
				},
			},
		},
		{
			name:    "map",
			val:     map[string]string{},
			wantErr: errs.ErrPointerOnly,
		},
		{
			name:    "slice",
			val:     []int{},
			wantErr: errs.ErrPointerOnly,
		},
		{
			name:    "basic type",
			val:     0,
			wantErr: errs.ErrPointerOnly,
		},

		// 标签相关测试用例
		{
			name: "column tag",
			val: func() any {
				// 我们把测试结构体定义在方法内部，防止被其它用例访问
				type ColumnTag struct {
					ID uint64 `orm:"column=id_123"`
				}
				return &ColumnTag{}
			}(),
			wantModel: &model{
				tableName: "column_tag",
				fieldMap: map[string]*field{
					"ID": {
						colName: "id_123",
					},
				},
			},
		},
		{
			// 如果用户设置了 column，但是传入一个空字符串，那么会用默认的名字
			name: "empty column",
			val: func() any {
				// 我们把测试结构体定义在方法内部，防止被其它用例访问
				type EmptyColumn struct {
					FirstName uint64 `orm:"column="`
				}
				return &EmptyColumn{}
			}(),
			wantModel: &model{
				tableName: "empty_column",
				fieldMap: map[string]*field{
					"FirstName": {
						colName: "first_name",
					},
				},
			},
		},
		{
			// 如果用户设置了 column，但是没有赋值
			name: "invalid Tag",
			val: func() any {
				type Invalid struct {
					FirstName string `orm:"column"`
				}
				return &Invalid{}
			}(),
			wantErr: errs.NewErrInvalidTag("column"),
		},
		{
			// 如果用户设置了一些奇奇怪怪的内容，这部分内容我们会忽略掉
			name: "ignore tag",
			val: func() any {
				// 我们把测试结构体定义在方法内部，防止被其它用例访问
				type IgnoreTag struct {
					FirstName uint64 `orm:"abc=abc"`
				}
				return &IgnoreTag{}
			}(),
			wantModel: &model{
				tableName: "ignore_tag",
				fieldMap: map[string]*field{
					"FirstName": {
						colName: "first_name",
					},
				},
			},
		},
	}

	r := &registry{}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			m, err := r.get(tc.val)
			assert.Equal(t, tc.wantErr, err)
			if err != nil {
				return
			}
			assert.Equal(t, tc.wantModel, m)
		})
	}
}
