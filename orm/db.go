package orm

type DB struct {
	r *registry
}

type Option func(*DB)

func NewDB(opts ...Option) (*DB, error) {
	db := &DB{
		r: &registry{},
	}

	for _, opt := range opts {
		opt(db)
	}
	return db, nil
}

// MustNewDB 创建一个 DB，如果失败则会 panic
// 我个人不太喜欢这种
func MustNewDB(opts ...Option) *DB {
	db, err := NewDB(opts...)
	if err != nil {
		panic(err)
	}
	return db
}

//// 按理说 NewSelector 之类的东西应该是定义在 DB 之上的
//// 但是因为泛型的限制不能采用这种方法
//func (d *DB) NewSelector[T any]() Selector[T] {
//	return &Selector[T]{
//		db: d,
//	}
//}
