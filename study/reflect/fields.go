package reflect

import (
	"errors"
	"reflect"
)

// IterateFields 遍历获取字段的类型和值
// 只有导出字段才能获取到值
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

// SetField 设置字段的值
// 只有canset == true 的字段才能设置值, 否则panic
// canset == true 等价于 导出字段 + 可寻址 (可寻址意味着传的是引用, 对于结构体类型变量而言,就是必须使用结构体指针， 那么结
// 构体的字段才是可以修改的)
// 详细参考可见: Go语言精进之路2  59.4节
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
