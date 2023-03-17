// Package orm 解析模型数据(元数据)
package orm

import (
	"geektime-go-study/orm/internal/errs"
	"geektime-go-study/orm/internal/util"
	"reflect"
)

type model struct {
	tableName string
	fieldMap  map[string]*field
}

type field struct {
	colName string
}

// ParseModel 解析模型数据 支持用户传入结构体指针/结构体
func ParseModel(entity any) (*model, error) {
	if entity == nil {
		return nil, errs.ErrPointerOnly
	}
	typ := reflect.TypeOf(entity)
	for typ.Kind() == reflect.Ptr {
		typ = typ.Elem()
	}

	if typ.Kind() != reflect.Struct {
		return nil, errs.ErrPointerOnly
	}

	numField := typ.NumField()

	model := &model{
		tableName: util.CamelToUnderline(typ.Name()),
		fieldMap:  make(map[string]*field, numField),
	}

	for i := 0; i < numField; i++ {
		name := typ.Field(i).Name
		model.fieldMap[name] = &field{colName: util.CamelToUnderline(name)}
	}

	return model, nil
}
