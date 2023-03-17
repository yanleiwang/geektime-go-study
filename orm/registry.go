package orm

import (
	"geektime-go-study/orm/internal/errs"
	"geektime-go-study/orm/internal/util"
	"reflect"
	"strings"
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
	typ := reflect.TypeOf(entity)
	if m, ok := r.models.Load(typ); ok {
		return m.(*model), nil
	}

	m, err := r.parseModel(entity)
	if err != nil {
		return nil, err
	}
	r.models.Store(typ, m)
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
	fds := make(map[string]*field, numField)

	for i := 0; i < numField; i++ {
		fdType := typ.Field(i)
		name := fdType.Name
		ormTags, err := r.parseTag(fdType.Tag)
		if err != nil {
			return nil, err
		}

		colName := ormTags[tagKeyColumn]
		if colName == "" {
			colName = util.CamelToUnderline(name)
		}

		fds[name] = &field{colName: colName}
	}

	var tableName string
	if v, ok := entity.(TableName); ok {
		tableName = v.TableName()
	}
	if tableName == "" {
		tableName = util.CamelToUnderline(typ.Name())
	}

	return &model{
		tableName: tableName,
		fieldMap:  fds,
	}, nil
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

func (r *registry) parseTag(tag reflect.StructTag) (map[string]string, error) {
	ormTag := tag.Get("orm")
	if ormTag == "" {
		return map[string]string{}, nil
	}

	pairs := strings.Split(ormTag, ",")
	res := make(map[string]string, len(pairs))
	for _, pair := range pairs {
		kv := strings.Split(pair, "=")
		if len(kv) != 2 {
			return nil, errs.NewErrInvalidTag(pair)
		}

		res[kv[0]] = kv[1]
	}
	return res, nil
}
