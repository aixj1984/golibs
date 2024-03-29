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
	"fmt"
	"reflect"
	"strings"

	"github.com/aixj1984/golibs/gorm-plus/constants"
)

// QueryCond 是每一次查询的条件对象
type QueryCond[T any] struct {
	selectColumns    []string
	omitColumns      []string
	distinctColumns  []string
	queryExpressions []any
	orderBuilder     strings.Builder
	groupBuilder     strings.Builder
	havingBuilder    strings.Builder
	havingArgs       []any
	queryArgs        []any
	last             any
	limit            *int
	offset           int
	updateMap        map[string]any
	// columnTypeMap    map[string]reflect.Type
}

func (q *QueryCond[T]) getSQLSegment() string {
	return ""
}

// NewQuery 构建查询条件
func NewQuery[T any]() (*QueryCond[T], *T) {
	q := &QueryCond[T]{}

	modelTypeStr := reflect.TypeOf((*T)(nil)).Elem().String()
	if model, ok := modelInstanceCache.Load(modelTypeStr); ok {
		m, isReal := model.(*T)
		if isReal {
			return q, m
		}
	}
	m := new(T)
	Cache(m)
	return q, m
}

// NewQueryModel 构建查询条件
func NewQueryModel[T any, R any]() (*QueryCond[T], *T, *R) {
	q := &QueryCond[T]{}
	var t *T
	var r *R
	entityTypeStr := reflect.TypeOf((*T)(nil)).Elem().String()
	if model, ok := modelInstanceCache.Load(entityTypeStr); ok {
		m, isReal := model.(*T)
		if isReal {
			t = m
		}
	}

	modelTypeStr := reflect.TypeOf((*R)(nil)).Elem().String()
	if model, ok := modelInstanceCache.Load(modelTypeStr); ok {
		m, isReal := model.(*R)
		if isReal {
			r = m
		}
	}

	if t == nil {
		t = new(T)
		Cache(t)
	}

	if r == nil {
		r = new(R)
		Cache(r)
	}

	return q, t, r
}

// Eq 等于 =
func (q *QueryCond[T]) Eq(column any, val any) *QueryCond[T] {
	q.addExpression(q.buildSQLSegment(column, constants.Eq, val)...)
	return q
}

// Ne 不等于 !=
func (q *QueryCond[T]) Ne(column any, val any) *QueryCond[T] {
	q.addExpression(q.buildSQLSegment(column, constants.Ne, val)...)
	return q
}

// Gt 大于 >
func (q *QueryCond[T]) Gt(column any, val any) *QueryCond[T] {
	q.addExpression(q.buildSQLSegment(column, constants.Gt, val)...)
	return q
}

// Ge 大于等于 >=
func (q *QueryCond[T]) Ge(column any, val any) *QueryCond[T] {
	q.addExpression(q.buildSQLSegment(column, constants.Ge, val)...)
	return q
}

// Lt 小于 <
func (q *QueryCond[T]) Lt(column any, val any) *QueryCond[T] {
	q.addExpression(q.buildSQLSegment(column, constants.Lt, val)...)
	return q
}

// Le 小于等于 <=
func (q *QueryCond[T]) Le(column any, val any) *QueryCond[T] {
	q.addExpression(q.buildSQLSegment(column, constants.Le, val)...)
	return q
}

// Like 模糊 LIKE '%值%'
func (q *QueryCond[T]) Like(column any, val any) *QueryCond[T] {
	s := fmt.Sprintf("%v", val)
	q.addExpression(q.buildSQLSegment(column, constants.Like, "%"+s+"%")...)
	return q
}

// NotLike 非模糊 NOT LIKE '%值%'
func (q *QueryCond[T]) NotLike(column any, val any) *QueryCond[T] {
	s := fmt.Sprintf("%v", val)
	q.addExpression(q.buildSQLSegment(column, constants.Not+" "+constants.Like, "%"+s+"%")...)
	return q
}

// LikeLeft 左模糊 LIKE '%值'
func (q *QueryCond[T]) LikeLeft(column any, val any) *QueryCond[T] {
	s := fmt.Sprintf("%v", val)
	q.addExpression(q.buildSQLSegment(column, constants.Like, "%"+s)...)
	return q
}

// NotLikeLeft 非左模糊 NOT LIKE '%值'
func (q *QueryCond[T]) NotLikeLeft(column any, val any) *QueryCond[T] {
	s := fmt.Sprintf("%v", val)
	q.addExpression(q.buildSQLSegment(column, constants.Not+" "+constants.Like, "%"+s)...)
	return q
}

// LikeRight 右模糊 LIKE '值%'
func (q *QueryCond[T]) LikeRight(column any, val any) *QueryCond[T] {
	s := fmt.Sprintf("%v", val)
	q.addExpression(q.buildSQLSegment(column, constants.Like, s+"%")...)
	return q
}

// NotLikeRight 非右模糊 NOT LIKE '值%'
func (q *QueryCond[T]) NotLikeRight(column any, val any) *QueryCond[T] {
	s := fmt.Sprintf("%v", val)
	q.addExpression(q.buildSQLSegment(column, constants.Not+" "+constants.Like, s+"%")...)
	return q
}

// IsNull 是否为空 字段 IS NULL
func (q *QueryCond[T]) IsNull(column any) *QueryCond[T] {
	q.addExpression(q.buildSQLSegment(column, constants.IsNull, nil)...)
	return q
}

// IsNotNull 是否非空 字段 IS NOT NULL
func (q *QueryCond[T]) IsNotNull(column any) *QueryCond[T] {
	q.addExpression(q.buildSQLSegment(column, constants.IsNotNull, nil)...)
	return q
}

// In 字段 IN (值1, 值2, ...)
func (q *QueryCond[T]) In(column any, val any) *QueryCond[T] {
	q.addExpression(q.buildSQLSegment(column, constants.In, val)...)
	return q
}

// NotIn 字段 NOT IN (值1, 值2, ...)
func (q *QueryCond[T]) NotIn(column any, val any) *QueryCond[T] {
	q.addExpression(q.buildSQLSegment(column, constants.Not+" "+constants.In, val)...)
	return q
}

// Between BETWEEN 值1 AND 值2
func (q *QueryCond[T]) Between(column any, start, end any) *QueryCond[T] {
	q.addExpression(q.buildSQLSegment(column, constants.Between, start, constants.And, end)...)
	return q
}

// NotBetween NOT BETWEEN 值1 AND 值2
func (q *QueryCond[T]) NotBetween(column any, start, end any) *QueryCond[T] {
	q.addExpression(q.buildSQLSegment(column, constants.Not+" "+constants.Between, start, constants.And, end)...)
	return q
}

// Distinct 去除重复字段值
func (q *QueryCond[T]) Distinct(columns ...any) *QueryCond[T] {
	for _, v := range columns {
		q.distinctColumns = append(q.distinctColumns, getColumnName(v))
	}
	return q
}

// Group 分组：GROUP BY 字段1,字段2
func (q *QueryCond[T]) Group(columns ...any) *QueryCond[T] {
	for _, v := range columns {
		columnName := getColumnName(v)
		if q.groupBuilder.Len() > 0 {
			q.groupBuilder.WriteString(constants.Comma) //nolint
		}
		q.groupBuilder.WriteString(columnName) //nolint
	}
	return q
}

// OrderByDesc 排序：ORDER BY 字段1,字段2 Desc
func (q *QueryCond[T]) OrderByDesc(columns ...any) *QueryCond[T] {
	var columnNames []string
	for _, v := range columns {
		columnName := getColumnName(v)
		columnNames = append(columnNames, columnName)
	}
	q.buildOrder(constants.Desc, columnNames...)
	return q
}

// OrderByAsc 排序：ORDER BY 字段1,字段2 ASC
func (q *QueryCond[T]) OrderByAsc(columns ...any) *QueryCond[T] {
	var columnNames []string
	for _, v := range columns {
		columnName := getColumnName(v)
		columnNames = append(columnNames, columnName)
	}
	q.buildOrder(constants.Asc, columnNames...)
	return q
}

// Having HAVING SQl语句
func (q *QueryCond[T]) Having(having string, args ...any) *QueryCond[T] {
	q.havingBuilder.WriteString(having) //nolint
	if len(args) == 1 {
		// 兼容function方法中in返回切片类型数据
		if anies, ok := args[0].([]any); ok {
			q.havingArgs = append(q.havingArgs, anies...)
			return q
		}
	}
	q.havingArgs = append(q.havingArgs, args...)
	return q
}

// And 拼接 AND
func (q *QueryCond[T]) And(fn ...func(q *QueryCond[T])) *QueryCond[T] {
	q.addExpression(&sqlKeyword{keyword: constants.And})
	if len(fn) > 0 {
		nestQuery := &QueryCond[T]{}
		fn[0](nestQuery)
		q.queryExpressions = append(q.queryExpressions, nestQuery)
		q.last = nestQuery
		return q
	}
	return q
}

// Or 拼接 OR
func (q *QueryCond[T]) Or(fn ...func(q *QueryCond[T])) *QueryCond[T] {
	q.addExpression(&sqlKeyword{keyword: constants.Or})
	if len(fn) > 0 {
		nestQuery := &QueryCond[T]{}
		fn[0](nestQuery)
		q.queryExpressions = append(q.queryExpressions, nestQuery)
		q.last = nestQuery
		return q
	}
	return q
}

// Select 查询字段
func (q *QueryCond[T]) Select(columns ...any) *QueryCond[T] {
	for _, v := range columns {
		columnName := getColumnName(v)
		q.selectColumns = append(q.selectColumns, columnName)
	}
	return q
}

// Omit 忽略字段
func (q *QueryCond[T]) Omit(columns ...any) *QueryCond[T] {
	for _, v := range columns {
		columnName := getColumnName(v)
		q.omitColumns = append(q.omitColumns, columnName)
	}
	return q
}

// Set 设置更新的字段
func (q *QueryCond[T]) Set(column any, val any) *QueryCond[T] {
	columnName := getColumnName(column)
	if q.updateMap == nil {
		q.updateMap = make(map[string]any)
	}
	q.updateMap[columnName] = val
	return q
}

func (q *QueryCond[T]) addExpression(sqlSegments ...SQLSegment) {
	if len(sqlSegments) == 1 {
		q.handleSingle(sqlSegments[0])
		return
	}

	// 如果符合条件，则添加and。
	q.addAndCondIfNeed()

	for _, sqlSegment := range sqlSegments {
		q.queryExpressions = append(q.queryExpressions, sqlSegment)
	}

	if len(sqlSegments) > 0 {
		q.last = sqlSegments[len(sqlSegments)-1]
	}
}

func (q *QueryCond[T]) addAndCondIfNeed() {
	// 1.如果上一个不是keyword，并且表达式切片不为空，则添加and。
	// 2.如果上一个值不是and或or,并且表达式切片不为空,则添加and。
	lastKeyword, isKeyword := q.last.(*sqlKeyword)
	if isNotKeyword(isKeyword, q.queryExpressions) || isLastNotAndOr(lastKeyword, isKeyword, q.queryExpressions) {
		sk := sqlKeyword{keyword: constants.And}
		q.queryExpressions = append(q.queryExpressions, &sk)
		q.last = &sk
	}
}

func (q *QueryCond[T]) handleSingle(sqlSegment SQLSegment) {
	// 如何是第一次设置，则不需要添加and(),or(),防止用户首次设置条件错误
	if len(q.queryExpressions) == 0 {
		return
	}

	// 防止用户重复设置and(),or()
	isRepeat := q.handelRepeat(sqlSegment)

	// 如果不是重复设置，则添加
	if !isRepeat {
		q.queryExpressions = append(q.queryExpressions, sqlSegment)
		q.last = sqlSegment
	}
}

func (q *QueryCond[T]) handelRepeat(sqlSegment SQLSegment) bool {
	currentKeyword, isCurrentKeyword := sqlSegment.(*sqlKeyword)
	lastKeyword, isLastKeyword := q.last.(*sqlKeyword)
	if isCurrentKeyword && isLastKeyword {
		// 如果上一次是and，这一次也是and，那么就不需要重复设置了
		isAnd := lastKeyword.keyword == constants.And
		if isAnd && currentKeyword.keyword == constants.And {
			return true
		}
		// 如果上一次是or，这一次也是or，那么就不需要重复设置了
		isOr := lastKeyword.keyword == constants.Or
		if isOr && currentKeyword.keyword == constants.Or {
			return true
		}
		// 如果上一次是and，这次是or，那么删除上一次的值，使用最新的值
		// 如果上一次是or，这次是and，那么删除上一次的值，使用最新的值
		q.queryExpressions = append(q.queryExpressions, q.queryExpressions[:len(q.queryExpressions)-1]...)
	}
	return false
}

func isNotKeyword(isKeyword bool, expressions []any) bool {
	return !isKeyword && len(expressions) > 0
}

func isLastNotAndOr(lastKeyword *sqlKeyword, isKeyword bool, expressions []any) bool {
	return isKeyword && lastKeyword.keyword != constants.And && lastKeyword.keyword != constants.Or && len(expressions) > 0
}

func (q *QueryCond[T]) buildSQLSegment(column any, condType string, values ...any) []SQLSegment {
	var sqlSegments []SQLSegment
	sqlSegments = append(sqlSegments, &columnPointer{column: column}, &sqlKeyword{keyword: condType})
	for _, val := range values {
		cv := columnValue{value: val}
		sqlSegments = append(sqlSegments, &cv)
	}
	return sqlSegments
}

func (q *QueryCond[T]) buildOrder(orderType string, columns ...string) {
	for _, v := range columns {
		if q.orderBuilder.Len() > 0 {
			q.orderBuilder.WriteString(constants.Comma) //nolint
		}
		q.orderBuilder.WriteString(v)         //nolint
		q.orderBuilder.WriteString(" ")       //nolint
		q.orderBuilder.WriteString(orderType) //nolint
	}
}

// AddAndStrCond 执行增加AND条件
func (q *QueryCond[T]) AddAndStrCond(cond string) *QueryCond[T] {
	if len(q.queryExpressions) > 0 {
		sk := sqlKeyword{keyword: constants.And}
		q.queryExpressions = append(q.queryExpressions, &sk)
	}
	condSk := sqlKeyword{keyword: cond}
	q.queryExpressions = append(q.queryExpressions, &condSk)
	q.last = &condSk
	return q
}

// AddOrStrCond 执行增加OR条件
func (q *QueryCond[T]) AddOrStrCond(cond string) *QueryCond[T] {
	if len(q.queryExpressions) > 0 {
		sk := sqlKeyword{keyword: constants.Or}
		q.queryExpressions = append(q.queryExpressions, &sk)
	}
	condSk := sqlKeyword{keyword: cond}
	q.queryExpressions = append(q.queryExpressions, &condSk)
	q.last = &condSk
	return q
}

// Case 根据条件，执行方法
func (q *QueryCond[T]) Case(isTrue bool, handleFunc func()) *QueryCond[T] {
	if isTrue {
		handleFunc()
	}
	return q
}

// Reset 重置查询条件
func (q *QueryCond[T]) Reset() *QueryCond[T] {
	q.selectColumns = q.selectColumns[:0]
	q.omitColumns = q.omitColumns[:0]
	q.distinctColumns = q.distinctColumns[:0]
	q.queryExpressions = q.queryExpressions[:0]
	q.orderBuilder.Reset()
	q.groupBuilder.Reset()
	q.havingBuilder.Reset()
	q.havingArgs = q.havingArgs[:0]
	q.queryArgs = q.queryArgs[:0]
	q.updateMap = nil

	return q
}
