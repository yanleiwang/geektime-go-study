// Package orm 解析模型数据(元数据)
package orm

type model struct {
	tableName string
	fieldMap  map[string]*field
}

type field struct {
	colName string
}
