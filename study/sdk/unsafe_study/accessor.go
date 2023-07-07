package unsafe_study

import (
	"errors"
	"reflect"
	"unsafe"
)

/*
通过  reflect + unsafe 的方法 读写字段
*/
type UnsafeAccessor struct {
	fields  map[string]FieldMeta
	address unsafe.Pointer // 结构体起始地址
}

type FieldMeta struct {
	Offset uintptr      // 字段偏移量
	typ    reflect.Type // 字段类型
}

func NewUnsafeAccessor(entity any) *UnsafeAccessor {
	typ := reflect.TypeOf(entity)
	typ = typ.Elem()
	numField := typ.NumField()
	fields := make(map[string]FieldMeta, numField)
	for i := 0; i < numField; i++ {
		fd := typ.Field(i)
		fields[fd.Name] = FieldMeta{
			Offset: fd.Offset,
			typ:    fd.Type,
		}
	}
	val := reflect.ValueOf(entity)
	return &UnsafeAccessor{
		fields:  fields,
		address: val.UnsafePointer(),
	}
}

// 读： *(*T)(ptr)， T 是目标类型，
// 如果类型不知道， 只能拿到反射的 Type，
// 那么可以用 reflect.NewAt(typ, ptr).Elem()。
// ptr 是字段偏移量：
// ptr = 结构体起始地址 + 字段偏移量
func (a *UnsafeAccessor) Field(field string) (any, error) {
	// 起始地址 + 字段偏移量
	fd, ok := a.fields[field]
	if !ok {
		return nil, errors.New("非法字段")
	}
	// 字段起始地址
	fdAddress := unsafe.Pointer(uintptr(a.address) + fd.Offset)
	// 如果知道类型，就这么读
	// return *(*int)(fdAddress), nil

	// 不知道确切类型
	return reflect.NewAt(fd.typ, fdAddress).Elem().Interface(), nil
}

// 写： *(*T)(ptr) = val， T 是目标类型。
// 如果类型不知道， 只能拿到反射的 Type，
// 那么可以用 reflect.NewAt(T, ptr).Elem().Set(reflect.ValueOf(val))。
// ptr 是字段偏移量：
// ptr = 结构体起始地址 + 字段偏移量
func (a *UnsafeAccessor) SetField(field string, val any) error {
	// 起始地址 + 字段偏移量
	fd, ok := a.fields[field]
	if !ok {
		return errors.New("非法字段")
	}
	// 字段起始地址
	fdAddress := unsafe.Pointer(uintptr(a.address) + fd.Offset)

	// 你知道确切类型就这么写
	// *(*int)(fdAddress) = val.(int)

	// 你不知道确切类型
	reflect.NewAt(fd.typ, fdAddress).Elem().Set(reflect.ValueOf(val))
	return nil
}
