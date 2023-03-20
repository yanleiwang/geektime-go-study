package model

import (
	"geektime-go-study/orm/internal/errs"
	"geektime-go-study/orm/internal/util"
	"reflect"
	"strings"
	"sync"
)

// Registry 元数据注册中心的抽象
// 允许用户显式地注册模型  所以 type registry struct的构造函数NewRegistry是导出的
type Registry interface {
	// Get 查找元数据
	Get(val any) (*Model, error)
	// Register 注册一个模型
	Register(val any, opts ...Option) (*Model, error)
}

// 这种包变量对测试不友好，缺乏隔离
// var defaultRegistry = &registry{
// 	models: make(map[reflect.Type]*Model, 16),
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

func NewRegistry() Registry {
	return &registry{}
}

func (r *registry) Get(entity any) (*Model, error) {
	typ := reflect.TypeOf(entity)
	m, ok := r.models.Load(typ)
	if ok {
		return m.(*Model), nil
	}

	return r.Register(entity)
}

func (r *registry) Register(val any, opts ...Option) (*Model, error) {
	m, err := r.parseModel(val)
	if err != nil {
		return nil, err
	}

	for _, opt := range opts {
		err = opt(m)
		if err != nil {
			return nil, err
		}
	}

	typ := reflect.TypeOf(val)
	r.models.Store(typ, m)
	return m, nil
}

// ParseModel 解析模型数据 支持用户传入结构体指针/结构体
// 支持从标签中提取自定义设置
// 标签形式 orm:"key1=value1,key2=value2"
func (r *registry) parseModel(entity any) (*Model, error) {
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
	fds := make(map[string]*Field, numField)
	cols := make(map[string]*Field, numField)

	for i := 0; i < numField; i++ {
		fdType := typ.Field(i)
		fdName := fdType.Name
		ormTags, err := r.parseTag(fdType.Tag)
		if err != nil {
			return nil, err
		}

		colName := ormTags[tagKeyColumn]
		if colName == "" {
			colName = util.CamelToUnderline(fdName)
		}

		f := &Field{
			ColName:   colName,
			FieldType: fdType.Type,
			FieldName: fdName,
			Offset:    fdType.Offset,
		}

		fds[fdName] = f
		cols[colName] = f
	}

	var tableName string
	if v, ok := entity.(TableName); ok {
		tableName = v.TableName()
	}
	if tableName == "" {
		tableName = util.CamelToUnderline(typ.Name())
	}

	return &Model{
		TableName: tableName,
		FieldMap:  fds,
		ColMap:    cols,
	}, nil
}

// 直接 map
// func (r *registry) get(val any) (*Model, error) {
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
// func (r *registry) get1(val any) (*Model, error) {
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
