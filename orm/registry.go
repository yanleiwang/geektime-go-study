package orm

import (
	"geektime-go-study/orm/internal/errs"
	"geektime-go-study/orm/internal/util"
	"reflect"
	"sync"
)

// 这种包变量对测试不友好，缺乏隔离
// var defaultRegistry = &registry{
// 	models: make(map[reflect.Type]*model, 16),
// }

type registry struct {
	// models key 是类型名
	// 这种定义方式是不行的
	// 1. 类型名冲突，例如都是 User，但是一个映射过去 buyer_t
	// 一个映射过去 seller_t
	// 2. 并发不安全
	// models map[string]*model

	// lock sync.RWMutex
	// models map[reflect.Type]*model
	models sync.Map
}

func (r *registry) get(entity any) (*model, error) {
	if m, ok := r.models.Load(reflect.TypeOf(entity)); ok {
		return m.(*model), nil
	}

	m, err := r.parseModel(entity)
	if err != nil {
		return nil, err
	}
	r.models.Store(reflect.TypeOf(entity), m)
	return m, nil
}

// ParseModel 解析模型数据 支持用户传入结构体指针/结构体
func (r *registry) parseModel(entity any) (*model, error) {
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

// 直接 map
// func (r *registry) get(val any) (*model, error) {
// 	typ := reflect.TypeOf(val)
// 	m, ok := r.models[typ]
// 	if !ok {
// 		var err error
// 		if m, err = r.parseModel(typ); err != nil {
// 			return nil, err
// 		}
// 	}
// 	r.models[typ] = m
// 	return m, nil
// }

// 使用读写锁的并发安全解决思路
// func (r *registry) get1(val any) (*model, error) {
// 	r.lock.RLock()
// 	typ := reflect.TypeOf(val)
// 	m, ok := r.models[typ]
// 	r.lock.RUnlock()
// 	if ok {
// 		return m, nil
// 	}
// 	r.lock.Lock()
// 	defer r.lock.Unlock()
// 	m, ok = r.models[typ]
// 	if ok {
// 		return m, nil
// 	}
// 	var err error
// 	if m, err = r.parseModel(typ); err != nil {
// 		return nil, err
// 	}
// 	r.models[typ] = m
// 	return m, nil
// }
