package orm

import (
	"database/sql"
	"geektime-go-study/orm/internal/errs"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestParseModel(t *testing.T) {
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
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			model, err := ParseModel(tc.val)
			assert.Equal(t, tc.wantErr, err)
			if err != nil {
				return
			}
			assert.Equal(t, tc.wantModel, model)
		})
	}

}
