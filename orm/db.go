package orm

import (
	"database/sql"
	"geektime-go-study/orm/internal/valuer"
	"geektime-go-study/orm/model"
)

type DB struct {
	r          model.Registry // 元数据注册中心
	db         *sql.DB
	valCreator valuer.Creator // 负责创建结构体的抽象(反射 or unsafe 实现, 默认unsafe实现)
}

type DBOption func(*DB)

func DBWithReflectValuer() DBOption {
	return func(db *DB) {
		db.valCreator = valuer.NewReflectValue
	}
}

func DBWithRegistry(r model.Registry) DBOption {
	return func(db *DB) {
		db.r = r
	}
}

func Open(driver string, dsn string, opts ...DBOption) (*DB, error) {
	db, err := sql.Open(driver, dsn)
	if err != nil {
		return nil, err
	}

	return OpenDB(db, opts...)
}

func OpenDB(db *sql.DB, opts ...DBOption) (*DB, error) {
	ret := &DB{
		r:          model.NewRegistry(),
		db:         db,
		valCreator: valuer.NewUnsafeValue,
	}

	for _, opt := range opts {
		opt(ret)
	}
	return ret, nil
}

// MustOpen 创建一个 DB，如果失败则会 panic
// 我个人不太喜欢这种
func MustOpen(driver string, dsn string, opts ...DBOption) *DB {
	ret, err := Open(driver, dsn, opts...)
	if err != nil {
		panic(err)
	}

	return ret
}

//// 按理说 NewSelector 之类的东西应该是定义在 DB 之上的
//// 但是因为泛型的限制不能采用这种方法
//func (d *DB) NewSelector[T any]() Selector[T] {
//	return &Selector[T]{
//		db: d,
//	}
//}
