// Package errs 中心式的 error 定义 方便后期维护 重构
package errs

import (
	"errors"
	"fmt"
)

var (
	ErrUnsupportedExpressionType = errors.New("orm: 不支持的表达式")
	// ErrPointerOnly 只支持一级指针作为输入
	// 看到这个 error 说明你输入了其它的东西
	// 我们并不希望用户能够直接使用 err == ErrPointerOnly
	// 所以放在我们的 internal 包里
	ErrPointerOnly            = errors.New("orm: 只支持一级指针作为输入，例如 *User")
	ErrNoRows                 = errors.New("orm: 未找到数据")
	ErrTooManyReturnedColumns = errors.New("orm: 过多列")
)

func NewErrUnsupportedExpressionType(exp any) error {
	return fmt.Errorf("%w %v", ErrUnsupportedExpressionType, exp)
}

// NewErrUnknownField 返回代表未知字段的错误
// 一般意味着你可能输入的是列名，或者输入了错误的字段名
func NewErrUnknownField(fd string) error {
	return fmt.Errorf("orm: 未知字段 %s", fd)
}

func NewErrInvalidTag(tag string) error {
	return fmt.Errorf("orm: 非法tag %s", tag)
}

func NewErrUnknownColumn(colName string) error {
	return fmt.Errorf("orm: 未知列 %s", colName)
}
