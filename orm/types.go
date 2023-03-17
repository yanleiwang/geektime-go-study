package orm

import (
	"context"
	"database/sql"
)

// Querier 使用泛型做类型约束
type Querier[T any] interface {
	Get(ctx context.Context) (*T, error)
	GetMulti(ctx context.Context) ([]*T, error)
}

type Executor interface {
	Exec(ctx context.Context) (sql.Result, error)
}

type Query struct {
	SQL  string
	Args []any
}

// QueryBuilder 作为构建 SQL 这一个单独步骤的顶级抽象
// 目的是把 构建sql 和 发起查询分开 这样程序逻辑清晰 方便测试
type QueryBuilder interface {
	Build() (*Query, error)
}
