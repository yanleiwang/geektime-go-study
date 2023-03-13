package reflect

import (
	"errors"
	"reflect"
)

func IterateFields(entity any) (map[string]any, error) {
	if entity == nil {
		return nil, errors.New("不支持 nil")
	}
	typ := reflect.TypeOf(entity)
	val := reflect.ValueOf(entity)

	for typ.Kind() == reflect.Ptr {
		// 拿到指针指向的对象
		typ = typ.Elem()
		val = val.Elem()
	}

	if !val.IsValid() {
		return nil, errors.New("不支持无效值")
	}

	if typ.Kind() != reflect.Struct {
		return nil, errors.New("不支持类型")
	}

	num := typ.NumField()
	res := make(map[string]any, num)
	for i := 0; i < num; i++ {
		// 字段的类型
		fieldType := typ.Field(i)
		// 字段的值
		fieldValue := val.Field(i)
		// 反射能够拿到私有字段的类型信息， 但是拿不到值,  所以 取其零值
		if fieldType.IsExported() {
			res[fieldType.Name] = fieldValue.Interface()
		} else {
			res[fieldType.Name] = reflect.Zero(fieldType.Type).Interface()
		}
	}
	return res, nil

}

func SetField(entity any, field string, newValue any) error {
	val := reflect.ValueOf(entity)
	for val.Type().Kind() == reflect.Ptr {
		val = val.Elem()
	}
	fieldVal := val.FieldByName(field)
	// 修改字段的值之前一定要先检查 CanSet
	if !fieldVal.CanSet() {
		return errors.New("不可修改字段")
	}
	fieldVal.Set(reflect.ValueOf(newValue))
	return nil

}
