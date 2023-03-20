package valuer

import (
	"database/sql"
	"geektime-go-study/orm/model"
)

// Valuer 是对结构体实例的内部抽象
type Valuer interface {
	// SetColumns 设置新值
	SetColumns(rows *sql.Rows) error
}

type Creator func(val any, meta *model.Model) Valuer

// ResultSetHandler 这是另外一种可行的设计方案
// type ResultSetHandler interface {
// 	// SetColumns 设置新值，column 是列名
// 	SetColumns(val any, rows *sql.Rows) error
// }
