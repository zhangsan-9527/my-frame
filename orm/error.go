package orm

import "my-frame/orm/internal/errs"

// ErrNoRows 通过这种形式将内部错误, 暴露在外面
var ErrNoRows = errs.ErrNoRows
