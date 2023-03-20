package model

import (
	"database/sql"
	"errors"
	"geektime-go-study/orm/internal/errs"
	"github.com/stretchr/testify/assert"
	"reflect"
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
		wantModel *Model
		wantErr   error
	}{
		{
			// 指针
			name: "pointer",
			val:  &TestModel{},
			wantModel: &Model{
				TableName: "test_model",
				FieldMap: map[string]*Field{
					"Id": {
						ColName:   "id",
						FieldType: reflect.TypeOf(int64(0)),
						FieldName: "Id",
						Offset:    0,
					},
					"FirstName": {
						ColName:   "first_name",
						FieldType: reflect.TypeOf(""),
						FieldName: "FirstName",
						Offset:    8,
					},
					"Age": {
						ColName:   "age",
						FieldType: reflect.TypeOf(int8(0)),
						FieldName: "Age",
						Offset:    24,
					},
					"LastName": {
						ColName:   "last_name",
						FieldType: reflect.TypeOf(&sql.NullString{}),
						FieldName: "LastName",
						Offset:    32,
					},
				},
				ColMap: map[string]*Field{
					"id": {
						ColName:   "id",
						FieldType: reflect.TypeOf(int64(0)),
						FieldName: "Id",
						Offset:    0,
					},
					"first_name": {
						ColName:   "first_name",
						FieldType: reflect.TypeOf(""),
						FieldName: "FirstName",
						Offset:    8,
					},
					"age": {
						ColName:   "age",
						FieldType: reflect.TypeOf(int8(0)),
						FieldName: "Age",
						Offset:    24,
					},
					"last_name": {
						ColName:   "last_name",
						FieldType: reflect.TypeOf(&sql.NullString{}),
						FieldName: "LastName",
						Offset:    32,
					},
				},
			},
		},
		{
			name:    "map",
			val:     map[string]string{},
			wantErr: errors.New("orm: 只支持一级指针作为输入，例如 *User"),
		},
		{
			name:    "slice",
			val:     []int{},
			wantErr: errors.New("orm: 只支持一级指针作为输入，例如 *User"),
		},
		{
			name:    "basic type",
			val:     0,
			wantErr: errors.New("orm: 只支持一级指针作为输入，例如 *User"),
		},

		// 标签相关测试用例
		{
			name: "column tag",
			val: func() any {
				// 我们把测试结构体定义在方法内部，防止被其它用例访问
				type ColumnTag struct {
					ID uint64 `orm:"column=id"`
				}
				return &ColumnTag{}
			}(),
			wantModel: &Model{
				TableName: "column_tag",
				FieldMap: map[string]*Field{
					"ID": {
						ColName:   "id",
						FieldType: reflect.TypeOf(uint64(0)),
						FieldName: "ID",
					},
				},
				ColMap: map[string]*Field{
					"id": {
						ColName:   "id",
						FieldType: reflect.TypeOf(uint64(0)),
						FieldName: "ID",
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
					FirstName string `orm:"column="`
				}
				return &EmptyColumn{}
			}(),
			wantModel: &Model{
				TableName: "empty_column",
				FieldMap: map[string]*Field{
					"FirstName": {
						ColName:   "first_name",
						FieldType: reflect.TypeOf(""),
						FieldName: "FirstName",
					},
				},
				ColMap: map[string]*Field{
					"first_name": {
						ColName:   "first_name",
						FieldType: reflect.TypeOf(""),
						FieldName: "FirstName",
					},
				},
			},
		},
		{
			// 如果用户设置了 column，但是没有赋值
			name: "invalid tag",
			val: func() any {
				// 我们把测试结构体定义在方法内部，防止被其它用例访问
				type InvalidTag struct {
					FirstName string `orm:"column"`
				}
				return &InvalidTag{}
			}(),
			wantErr: errs.NewErrInvalidTag("column"),
		},
		{
			// 如果用户设置了一些奇奇怪怪的内容，这部分内容我们会忽略掉
			name: "ignore tag",
			val: func() any {
				// 我们把测试结构体定义在方法内部，防止被其它用例访问
				type IgnoreTag struct {
					FirstName string `orm:"abc=abc"`
				}
				return &IgnoreTag{}
			}(),
			wantModel: &Model{
				TableName: "ignore_tag",
				FieldMap: map[string]*Field{
					"FirstName": {
						ColName:   "first_name",
						FieldType: reflect.TypeOf(""),
						FieldName: "FirstName",
					},
				},
				ColMap: map[string]*Field{
					"first_name": {
						ColName:   "first_name",
						FieldType: reflect.TypeOf(""),
						FieldName: "FirstName",
					},
				},
			},
		},

		// 利用接口自定义模型信息
		{
			name: "table name",
			val:  &CustomTableName{},
			wantModel: &Model{
				TableName: "custom_table_name_t",
				FieldMap: map[string]*Field{
					"Name": {
						ColName:   "name",
						FieldName: "Name",
						FieldType: reflect.TypeOf(""),
					},
				},
				ColMap: map[string]*Field{
					"name": {
						ColName:   "name",
						FieldName: "Name",
						FieldType: reflect.TypeOf(""),
					},
				},
			},
		},
		{
			name: "table name ptr",
			val:  &CustomTableNamePtr{},
			wantModel: &Model{
				TableName: "custom_table_name_ptr_t",
				FieldMap: map[string]*Field{
					"Name": {
						ColName:   "name",
						FieldName: "Name",
						FieldType: reflect.TypeOf(""),
					},
				},
				ColMap: map[string]*Field{
					"name": {
						ColName:   "name",
						FieldName: "Name",
						FieldType: reflect.TypeOf(""),
					},
				},
			},
		},
		{
			name: "empty table name",
			val:  &EmptyTableName{},
			wantModel: &Model{
				TableName: "empty_table_name",
				FieldMap: map[string]*Field{
					"Name": {
						ColName:   "name",
						FieldName: "Name",
						FieldType: reflect.TypeOf(""),
					},
				},
				ColMap: map[string]*Field{
					"name": {
						ColName:   "name",
						FieldName: "Name",
						FieldType: reflect.TypeOf(""),
					},
				},
			},
		},
	}

	r := NewRegistry()
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			m, err := r.Get(tc.val)
			assert.Equal(t, tc.wantErr, err)
			if err != nil {
				return
			}
			assert.Equal(t, tc.wantModel, m)
		})
	}
}

type CustomTableName struct {
	Name string
}

func (c CustomTableName) TableName() string {
	return "custom_table_name_t"
}

type CustomTableNamePtr struct {
	Name string
}

func (c *CustomTableNamePtr) TableName() string {
	return "custom_table_name_ptr_t"
}

type EmptyTableName struct {
	Name string
}

func (c *EmptyTableName) TableName() string {
	return ""
}
