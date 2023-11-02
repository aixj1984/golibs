package gorm

import (
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func Expr(expr string, args ...interface{}) clause.Expr {
	return gorm.Expr(expr, args...)
}
