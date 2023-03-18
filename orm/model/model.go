// Package model 解析模型数据(元数据)
package model

import "geektime-go-study/orm/internal/errs"

// Model 元数据
// Model是导出的原因是我们暴露了Register
type Model struct {
	TableName string
	FieldMap  map[string]*Field
}

type Field struct {
	ColName string
}

type ModOption func(model *Model) error

func WithTableName(name string) ModOption {
	return func(m *Model) error {
		m.TableName = name
		return nil
	}
}

func WithColumnName(field string, columnName string) ModOption {
	return func(m *Model) error {
		fd, ok := m.FieldMap[field]
		if !ok {
			return errs.NewErrUnknownField(field)
		}
		// 注意，这里我们根本没有检测 colName 会不会是空字符串
		// 因为正常情况下，用户都不会写错
		// 即便写错了，也很容易在测试中发现
		fd.ColName = columnName
		return nil
	}
}

// 我们支持的全部标签上的 key 都放在这里
// 方便用户查找，和我们后期维护
const (
	tagKeyColumn = "column"
)

// 用户自定义一些模型信息的接口，集中放在这里
// 方便用户查找和我们后期维护

// TableName 用户实现这个接口来返回自定义的表名
type TableName interface {
	TableName() string
}
