package errs

import (
	"errors"
	"fmt"
)

var (
	ErrPointerOnly = errors.New("orm: 只支持指向结构体的一级指针")

	//errUnsupportedExpression = errors.New("orm: 不支持的表达形式")
)

//func NewErrUnsupportedExpressionV1(expr any) error {
//	return fmt.Errorf("%w %v", errUnsupportedExpression, expr)
//}

// @ErrUnsupportedExpression 40001 原因是你输入了乱七八糟的类型
// 解决方案: 使用正确的类型
func NewErrUnsupportedExpression(expr any) error {
	return fmt.Errorf("orm: 不支持的表达形式 %v", expr)
}

func NewErrUnknownField(name any) error {
	return fmt.Errorf("orm: 未知字段 %v", name)
}

func NewErrInvalidTagContent(pair string) error {
	return fmt.Errorf("orm: 非法标签 %v", pair)
}
