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
// 查询条件可以看做是 Left Op Right 的模式：
// • 基本比较符：Left 是列名，Op 是各个比较符号，右
// 边是表达式，常见的是一个值
// • Not：左边缺省，只剩下 Op Right，如 NOT (id = ?)
// • And、Or：左边右边都是一个 Predicate
// 所以一个predicate left可以是predicate/列名/nil  right可以是值/predicate
// 我们把接口Expression 作为predicate/值/列名 的共同接口
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
