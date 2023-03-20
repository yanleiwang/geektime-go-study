package valuer

import (
	"database/sql"
	"geektime-go-study/orm/internal/errs"
	"geektime-go-study/orm/model"
	"reflect"
)

type reflectValue struct {
	val  reflect.Value
	meta *model.Model
}

// 确保 Creator 修改的时候, 能够得到提示
var _ Creator = NewReflectValue

func NewReflectValue(val any, meta *model.Model) Valuer {
	return &reflectValue{
		val:  reflect.ValueOf(val),
		meta: meta,
	}
}

func (r *reflectValue) SetColumns(rows *sql.Rows) error {
	// step 1 拿到结果集的列名
	colNames, err := rows.Columns()
	if err != nil {
		return err
	}

	if len(colNames) > len(r.meta.FieldMap) {
		return errs.ErrTooManyReturnedColumns
	}

	// step 2 根据列名 找到对应字段的元数据, 并设置字段type
	colVals := make([]any, 0, len(colNames))
	for _, colName := range colNames {
		cm, ok := r.meta.ColMap[colName]
		if !ok {
			return errs.NewErrUnknownColumn(colName)
		}

		colVal := reflect.New(cm.FieldType).Interface() // colVal 实质是指针
		colVals = append(colVals, colVal)
	}

	// step 3 拿到结果集的值
	err = rows.Scan(colVals...)
	if err != nil {
		return err
	}

	// step 4 把结果写入到 val中
	for i, colName := range colNames {
		cm := r.meta.ColMap[colName]
		fd := r.val.Elem().FieldByName(cm.FieldName)
		fd.Set(reflect.ValueOf(colVals[i]).Elem())
	}
	return nil
}
