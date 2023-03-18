package orm

import (
	"geektime-go-study/orm/internal/errs"
	"geektime-go-study/orm/model"
	"strings"
)

// Selector 使用泛型做类型约束
type Selector[T any] struct {
	tbl   string
	where []Predicate
	sb    strings.Builder
	args  []any
	m     *model.Model
	db    *DB
}

func (s *Selector[T]) Build() (*Query, error) {
	s.sb.Reset()
	s.sb.WriteString("SELECT * FROM ")
	// 决策：如果用户指定了表名，就直接使用，不会使用反引号；否则使用反引号括起来。
	var (
		t   T
		err error
	)
	s.m, err = s.db.r.Get(&t)
	if err != nil {
		return nil, err
	}
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
		field, ok := s.m.FieldMap[exp.name]
		if !ok {
			return errs.NewErrUnknownField(exp.name)
		}
		s.sb.WriteByte('`')
		s.sb.WriteString(field.ColName)
		s.sb.WriteByte('`')
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

func NewSelector[T any](db *DB) *Selector[T] {
	return &Selector[T]{
		db: db,
	}
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
