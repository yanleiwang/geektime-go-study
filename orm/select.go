package orm

import (
	"context"
	"geektime-go-study/orm/internal/errs"
	"geektime-go-study/orm/model"
	"strings"
)

// Selector 使用泛型做类型约束
type Selector[T any] struct {
	tbl     string
	where   []Predicate
	sb      strings.Builder
	args    []any
	m       *model.Model
	db      *DB
	columns []Selectable
}

func NewSelector[T any](db *DB) *Selector[T] {
	return &Selector[T]{
		db: db,
	}
}

func (s *Selector[T]) Get(ctx context.Context) (*T, error) {
	// step 1 构建sql
	query, err := s.Build()
	if err != nil {
		return nil, err
	}

	// step 2 发起查询
	// s.db 是我们定义的 DB
	// s.db.db 则是 sql.DB
	// 使用 QueryContext，从而和 GetMulti 能够复用处理结果集的代码
	rows, err := s.db.db.QueryContext(ctx, query.SQL, query.Args...)
	if err != nil {
		return nil, err
	}

	// 没有数据的话, 返回error 跟sql包语义一致
	if !rows.Next() {
		return nil, ErrNoRows
	}

	// step 3 结果集转为对象
	// step 3.1 创建对象
	val := new(T)
	// step 3.2 获取元数据
	meta, err := s.db.r.Get(val)
	if err != nil {
		return nil, err
	}
	// step 3.3 创建转换对象
	creator := s.db.valCreator(val, meta)
	// step 3.4 设置值
	err = creator.SetColumns(rows)
	if err != nil {
		return nil, err
	}
	return val, err

}

func (s *Selector[T]) GetMulti(ctx context.Context) ([]*T, error) {
	//TODO implement me
	panic("implement me")
}

func (s *Selector[T]) Build() (*Query, error) {
	s.sb.Reset()
	// 决策：如果用户指定了表名，就直接使用，不会使用反引号；否则使用反引号括起来。
	var (
		t   T
		err error
	)
	s.m, err = s.db.r.Get(&t)
	if err != nil {
		return nil, err
	}
	s.sb.WriteString("SELECT ")

	if len(s.columns) == 0 {
		s.sb.WriteString("*")
	} else {
		for i, c := range s.columns {
			if i != 0 {
				s.sb.WriteByte(',')
			}
			err := s.buildColumn(c.(Column))
			if err != nil {
				return nil, err
			}
		}
	}

	s.sb.WriteString(" FROM ")

	if s.tbl == "" {
		s.sb.WriteByte('`')
		s.sb.WriteString(s.m.TableName)
		s.sb.WriteByte('`')
	} else {
		s.sb.WriteString(s.tbl)
	}

	if len(s.where) > 0 {
		s.sb.WriteString(" WHERE ")
		p := s.where[0]
		// 用And 合并多个Predicate
		for i := 1; i < len(s.where); i++ {
			p = p.And(s.where[i])
		}

		if err := s.buildExpression(p); err != nil {
			return nil, err
		}

	}

	s.sb.WriteString(";")
	return &Query{
		SQL:  s.sb.String(),
		Args: s.args,
	}, nil
}

func (s *Selector[T]) buildExpression(e Expression) error {
	switch exp := e.(type) {
	case Predicate:
		_, lp := exp.left.(Predicate)
		if lp {
			s.sb.WriteByte('(')
		}
		if err := s.buildExpression(exp.left); err != nil {
			return err
		}
		if lp {
			s.sb.WriteByte(')')
		}

		s.sb.WriteByte(' ')
		s.sb.WriteString(exp.op.String())
		s.sb.WriteByte(' ')
		_, rp := exp.right.(Predicate)
		if rp {
			s.sb.WriteByte('(')
		}
		if err := s.buildExpression(exp.right); err != nil {
			return err
		}
		if rp {
			s.sb.WriteByte(')')
		}
	case Column:
		return s.buildColumn(exp)
	case value:
		s.sb.WriteByte('?')
		s.args = append(s.args, exp.val)
	case nil:
		return nil
	default:
		return errs.NewErrUnsupportedExpressionType(exp)
	}
	return nil
}

func (s *Selector[T]) buildColumn(c Column) error {
	field, ok := s.m.FieldMap[c.name]
	if !ok {
		return errs.NewErrUnknownField(c.name)
	}
	s.sb.WriteByte('`')
	s.sb.WriteString(field.ColName)
	s.sb.WriteByte('`')
	return nil
}

// From 考虑 FROM，可行的思路是:
// • Selector 本身有泛型参数，我们用泛型的类型名字作为表名
// • 加入一个 From 方法：如果用户调用了这个方法，那么我们就用这 个方法的参数来作为表名
func (s *Selector[T]) From(tbl string) *Selector[T] {
	s.tbl = tbl
	return s
}

// Where 设计1 接收字符串和参数作为输入
// 好处: 简单, 灵活
// 缺乏校验，用户容易写错，例如写错字段名，漏了括号等
// func (s *Selector[T]) Where(where string, args ...any)  {
//
// }
// 所以接收结构化的 Predicate 作为输入
func (s *Selector[T]) Where(ps ...Predicate) *Selector[T] {
	s.where = ps
	return s
}

func (s *Selector[T]) Select(cols ...Selectable) *Selector[T] {
	s.columns = cols
	return s
}

// Selectable 标记接口, 可以作为select xxx 里面 的xxx
// 有 Column
type Selectable interface {
	selectable()
}
