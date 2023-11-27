/*
 * Licensed to the AcmeStack under one or more contributor license
 * agreements. See the NOTICE file distributed with this work for
 * additional information regarding copyright ownership.
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *   http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package gplus

import (
	"strings"

	"github.com/aixj1984/golibs/gorm-plus/constants"
)

// Function 对数据库操作符的一个封装
type Function struct {
	funStr string
}

func (f *Function) As(asName any) string { //nolint
	return f.funStr + " " + constants.As + " " + getColumnName(asName)
}

func (f *Function) Eq(value int64) (string, int64) { //nolint
	return buildFunStr(f.funStr, constants.Eq, value)
}

func (f *Function) Ne(value int64) (string, int64) { //nolint
	return buildFunStr(f.funStr, constants.Ne, value)
}

func (f *Function) Gt(value int64) (string, int64) { //nolint
	return buildFunStr(f.funStr, constants.Gt, value)
}

func (f *Function) Ge(value int64) (string, int64) { //nolint
	return buildFunStr(f.funStr, constants.Ge, value)
}

func (f *Function) Lt(value int64) (string, int64) { //nolint
	return buildFunStr(f.funStr, constants.Lt, value)
}

func (f *Function) Le(value int64) (string, int64) { //nolint
	return buildFunStr(f.funStr, constants.Le, value)
}

func (f *Function) In(values ...any) (string, []any) { //nolint
	// 构建占位符
	placeholder := buildPlaceholder(values)
	return f.funStr + " " + constants.In + placeholder.String(), values
}

func (f *Function) NotIn(values ...any) (string, []any) { //nolint
	// 构建占位符
	placeholder := buildPlaceholder(values)
	return f.funStr + " " + constants.Not + " " + constants.In + placeholder.String(), values
}

func (f *Function) Between(start int64, end int64) (string, int64, int64) { //nolint
	return f.funStr + " " + constants.Between + " ? " + constants.And + " ?", start, end //nolint
}

func (f *Function) NotBetween(start int64, end int64) (string, int64, int64) { //nolint
	return f.funStr + " " + constants.Not + " " + constants.Between + " ? " + constants.And + " ?", start, end //nolint
}

func Sum(columnName any) *Function { //nolint
	return &Function{funStr: addBracket(constants.SUM, getColumnName(columnName))}
}

func Avg(columnName any) *Function { //nolint
	return &Function{funStr: addBracket(constants.AVG, getColumnName(columnName))}
}

func Max(columnName any) *Function { //nolint
	return &Function{funStr: addBracket(constants.MAX, getColumnName(columnName))}
}

func Min(columnName any) *Function { //nolint
	return &Function{funStr: addBracket(constants.MIN, getColumnName(columnName))}
}

func Count(columnName any) *Function { //nolint
	return &Function{funStr: addBracket(constants.COUNT, getColumnName(columnName))}
}

func As(columnName any, asName any) string { //nolint
	return getColumnName(columnName) + " " + constants.As + " " + getColumnName(asName)
}

func addBracket(function string, columnNameStr string) string { //nolint
	return function + constants.LeftBracket + columnNameStr + constants.RightBracket
}

func buildFunStr(funcStr string, typeStr string, value int64) (string, int64) { //nolint
	return funcStr + " " + typeStr + " ?", value
}

func buildPlaceholder(values []any) strings.Builder {
	var placeholder strings.Builder
	placeholder.WriteString(constants.LeftBracket) //nolint
	for i := 0; i < len(values); i++ {
		if i == len(values)-1 {
			placeholder.WriteString("?")                    //nolint
			placeholder.WriteString(constants.RightBracket) //nolint
			break
		}
		placeholder.WriteString("?") //nolint
		placeholder.WriteString(",") //nolint
	}
	return placeholder
}
