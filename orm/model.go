// Package orm 解析模型数据(元数据)
package orm

type model struct {
	tableName string
	fieldMap  map[string]*field
}

type field struct {
	colName string
}

// 我们支持的全部标签上的 key 都放在这里
// 方便用户查找，和我们后期维护
const (
	tagKeyColumn = "column"
)
