package orm

type value struct {
	val any
}

func (value) expr() {}

func valueOf(val any) value {
	return value{
		val: val,
	}
}

// Column 代表 某个列名
type Column struct {
	name string
}

// 实现Expression 标记接口
func (Column) expr() {}

// 标记 Selectable
func (Column) selectable() {}

func C(name string) Column {
	return Column{name: name}
}

// EQ 例如 C("id").EQ(12)
func (c Column) EQ(arg any) Predicate {
	return Predicate{
		left:  c,
		op:    opEQ,
		right: exprOf(arg), // 因为arg可能是Expression 所以在exprOf() 中进行判断
	}
}

func (c Column) LT(arg any) Predicate {
	return Predicate{
		left:  c,
		op:    opLT,
		right: exprOf(arg),
	}
}

func (c Column) GT(arg any) Predicate {
	return Predicate{
		left:  c,
		op:    opGT,
		right: exprOf(arg),
	}
}
