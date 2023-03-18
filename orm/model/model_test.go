package model

import (
	"database/sql"
	"geektime-go-study/orm/internal/errs"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestWithTableName(t *testing.T) {
	type TestModel struct {
		Id int64
		// ""
		FirstName string
		Age       int8
		LastName  *sql.NullString
	}

	testCases := []struct {
		name          string
		val           any
		opt           ModOption
		wantTableName string
		wantErr       error
	}{
		{
			// 我们没有对空字符串进行校验
			name:          "empty string",
			val:           &TestModel{},
			opt:           WithTableName(""),
			wantTableName: "",
		},
		{
			name:          "table name",
			val:           &TestModel{},
			opt:           WithTableName("test_model_t"),
			wantTableName: "test_model_t",
		},
	}

	r := NewRegistry()
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			m, err := r.Register(tc.val, tc.opt)
			assert.Equal(t, tc.wantErr, err)
			if err != nil {
				return
			}
			assert.Equal(t, tc.wantTableName, m.TableName)
		})
	}

}

func TestModelWithColumnName(t *testing.T) {
	type TestModel struct {
		Id int64
		// ""
		FirstName string
		Age       int8
		LastName  *sql.NullString
	}
	testCases := []struct {
		name        string
		val         any
		opt         ModOption
		field       string
		wantColName string
		wantErr     error
	}{
		{
			name:        "new name",
			val:         &TestModel{},
			opt:         WithColumnName("FirstName", "first_name_new"),
			field:       "FirstName",
			wantColName: "first_name_new",
		},
		{
			name:        "empty new name",
			val:         &TestModel{},
			opt:         WithColumnName("FirstName", ""),
			field:       "FirstName",
			wantColName: "",
		},
		{
			// 不存在的字段
			name:    "invalid field name",
			val:     &TestModel{},
			opt:     WithColumnName("FirstNameXXX", "first_name"),
			field:   "FirstNameXXX",
			wantErr: errs.NewErrUnknownField("FirstNameXXX"),
		},
	}

	r := NewRegistry()
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			m, err := r.Register(tc.val, tc.opt)
			assert.Equal(t, tc.wantErr, err)
			if err != nil {
				return
			}
			fd := m.FieldMap[tc.field]
			assert.Equal(t, tc.wantColName, fd.ColName)
		})
	}
}
