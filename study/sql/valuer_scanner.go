package sql

import (
	"database/sql"
	"database/sql/driver"
	"encoding/json"
	"fmt"
)

/*
SQL 默认支持的类型就是基础类型， []byte 和string，
该如何自定义类型？ 比如要支持 json 类型， 我该如何处理？
答案是实现以下两个接口:
• driver.Valuer： 读取， 实现该接口的类型可以作为查询参数使用(Go类型到数据库类型）
• sql.Scanner： 写入， 实现该接口的类型可以用于 Scan 方法（数据库类型到Go 类型）
自定义类型一般是实现这两个接口。
下面举例支持Json类型
*/

// JsonColumn 代表存储字段的 json 类型
// 主要用于没有提供默认 json 类型的数据库
// T 可以是结构体，也可以是切片或者 map
// 一切可以被 json 库所处理的类型都能被用作 T
type JsonColumn[T any] struct {
	Val   T
	Valid bool // Valid 用于解决null的问题
}

// Value 返回一个 json 串。类型是 []byte
// 实现 driver.Valuer 接口 , 用于支持向query/exec系列方法传递参数
func (j JsonColumn[T]) Value() (driver.Value, error) {
	if !j.Valid {
		return nil, nil
	}
	return json.Marshal(j.Val)
}

// Scan 将 src 转化为对象
// 因为是json字符串转 结构体,
// 所以src 的类型必须是 []byte, *[]byte, string, sql.RawBytes, *sql.RawBytes 之一
// 实现 sql.Scanner 接口 , 用于 Scan 方法
func (j *JsonColumn[T]) Scan(src any) error {
	var bs []byte
	switch val := src.(type) {
	case []byte:
		bs = val
	case *[]byte:
		bs = *val
	case string:
		bs = []byte(val)
	case sql.RawBytes:
		bs = val
	case *sql.RawBytes:
		bs = *val
	default:
		return fmt.Errorf("ekit：JsonColumn.Scan 不支持 src 类型 %v", src)
	}

	if err := json.Unmarshal(bs, &j.Val); err != nil {
		return err
	}
	j.Valid = true
	return nil
}
