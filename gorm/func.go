package gorm

import (
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// Expr 是gorm.Expr的透传
func Expr(expr string, args ...interface{}) clause.Expr {
	return gorm.Expr(expr, args...)
}
