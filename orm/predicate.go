package orm

type op string

const (
	opEQ  = "="
	opLT  = "<"
	opGT  = ">"
	opAND = "AND"
	opOR  = "OR"
	opNOT = "NOT"
)

func (o op) String() string {
	return string(o)
}

// Expression 代表语句，或者语句的部分
// 暂时没想好怎么设计方法，所以直接做成标记接口
type Expression interface {
	expr()
}

func exprOf(arg any) Expression {
	switch val := arg.(type) {
	case Expression:
		return val
	default:
		return valueOf(val)
	}
}

// Predicate 代表一个查询条件
// Predicate 可以通过和 Predicate 组合构成复杂的查询条件
type Predicate struct {
	left  Expression
	op    op
	right Expression
}

func (Predicate) expr() {}

func Not(p Predicate) Predicate {
	return Predicate{
		op:    opNOT,
		right: p,
	}
}

func (left Predicate) And(right Predicate) Predicate {
	return Predicate{
		left:  left,
		op:    opAND,
		right: right,
	}
}

func (left Predicate) Or(right Predicate) Predicate {
	return Predicate{
		left:  left,
		op:    opOR,
		right: right,
	}
}
