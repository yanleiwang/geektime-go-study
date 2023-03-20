package valuer

import (
	"database/sql"
	"geektime-go-study/orm/internal/errs"
	"geektime-go-study/orm/model"
	"reflect"
	"unsafe"
)

type unsafeValue struct {
	addr unsafe.Pointer
	meta *model.Model
}

var _ Creator = NewUnsafeValue

func NewUnsafeValue(val interface{}, meta *model.Model) Valuer {
	return &unsafeValue{
		addr: unsafe.Pointer(reflect.ValueOf(val).Pointer()),
		meta: meta,
	}
}

func (u *unsafeValue) SetColumns(rows *sql.Rows) error {
	cs, err := rows.Columns()
	if err != nil {
		return err
	}
	if len(cs) > len(u.meta.ColMap) {
		return errs.ErrTooManyReturnedColumns
	}

	colValues := make([]any, len(cs))
	for i, c := range cs {
		cm, ok := u.meta.ColMap[c]
		if !ok {
			return errs.NewErrUnknownColumn(c)
		}
		ptr := unsafe.Pointer(uintptr(u.addr) + cm.Offset)
		val := reflect.NewAt(cm.FieldType, ptr)
		colValues[i] = val.Interface()
	}

	return rows.Scan(colValues...)

}
