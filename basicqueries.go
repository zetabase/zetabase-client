package zetabase

import (
	"fmt"
	"github.com/zetabase/zetabase-client/zbprotocol"
)

// Nike: JUST DO IT

// Types for query "DSL"

type SubQueryConvertible interface {
	ToSubQuery(tblOwnerId, tblId string) *zbprotocol.TableSubQuery
}

type QueryAnd struct {
	Left  SubQueryConvertible
	Right SubQueryConvertible
}

type QueryOr struct {
	Left  SubQueryConvertible
	Right SubQueryConvertible
}

type QueryNotEquals struct {
	Field     string
	CompValue interface{}
}

type QueryEquals struct {
	Field     string
	CompValue interface{}
}

type QueryGreaterThan struct {
	Field     string
	CompValue interface{}
}

type QueryGreaterThanEqual struct {
	Field     string
	CompValue interface{}
}

type QueryLessThan struct {
	Field     string
	CompValue interface{}
}

type QueryLessThanEqual struct {
	Field     string
	CompValue interface{}
}

type QueryTextSearch struct {
	Field     string
	CompValue interface{}
}

func (q *QueryAnd) ToSubQuery(tblOwnerId, tblId string) *zbprotocol.TableSubQuery {
	return &zbprotocol.TableSubQuery{
		IsCompound:       true,
		CompoundOperator: zbprotocol.QueryLogicalOperator_LOGICAL_AND,
		CompoundLeft:     q.Left.ToSubQuery(tblOwnerId, tblId),
		CompoundRight:    q.Right.ToSubQuery(tblOwnerId, tblId),
		Comparison:       nil,
	}
}

func (q *QueryOr) ToSubQuery(tblOwnerId, tblId string) *zbprotocol.TableSubQuery {
	return &zbprotocol.TableSubQuery{
		IsCompound:       true,
		CompoundOperator: zbprotocol.QueryLogicalOperator_LOGICAL_OR,
		CompoundLeft:     q.Left.ToSubQuery(tblOwnerId, tblId),
		CompoundRight:    q.Right.ToSubQuery(tblOwnerId, tblId),
		Comparison:       nil,
	}
}


func QAnd(l, r SubQueryConvertible) *QueryAnd {
	return &QueryAnd{
		Left:  l,
		Right: r,
	}
}

func QOr(l, r SubQueryConvertible) *QueryOr {
	return &QueryOr{
		Left:  l,
		Right: r,
	}
}

func (q *QueryGreaterThanEqual) ToSubQuery(a, b string) *zbprotocol.TableSubQuery {
	valu := q.CompValue
	field := q.Field
	valuStr, qOrder := queryObjectTypify(valu)
	return &zbprotocol.TableSubQuery{
		IsCompound:       false,
		CompoundOperator: 0,
		CompoundLeft:     nil,
		CompoundRight:    nil,
		Comparison: &zbprotocol.TableSubqueryComparison{
			Op:       zbprotocol.QueryOperator_GREATER_THAN_EQ,
			Field:    field,
			Value:    valuStr,
			Ordering: qOrder,
		},
	}
}

func (q *QueryLessThanEqual) ToSubQuery(a, b string) *zbprotocol.TableSubQuery {
	valu := q.CompValue
	field := q.Field
	valuStr, qOrder := queryObjectTypify(valu)
	return &zbprotocol.TableSubQuery{
		IsCompound:       false,
		CompoundOperator: 0,
		CompoundLeft:     nil,
		CompoundRight:    nil,
		Comparison: &zbprotocol.TableSubqueryComparison{
			Op:       zbprotocol.QueryOperator_LESS_THAN_EQ,
			Field:    field,
			Value:    valuStr,
			Ordering: qOrder,
		},
	}
}

func (q *QueryLessThan) ToSubQuery(a, b string) *zbprotocol.TableSubQuery {
	valu := q.CompValue
	field := q.Field
	valuStr, qOrder := queryObjectTypify(valu)
	return &zbprotocol.TableSubQuery{
		IsCompound:       false,
		CompoundOperator: 0,
		CompoundLeft:     nil,
		CompoundRight:    nil,
		Comparison: &zbprotocol.TableSubqueryComparison{
			Op:       zbprotocol.QueryOperator_LESS_THAN,
			Field:    field,
			Value:    valuStr,
			Ordering: qOrder,
		},
	}
}


func (q *QueryGreaterThan) ToSubQuery(a, b string) *zbprotocol.TableSubQuery {
	valu := q.CompValue
	field := q.Field
	valuStr, qOrder := queryObjectTypify(valu)
	return &zbprotocol.TableSubQuery{
		IsCompound:       false,
		CompoundOperator: 0,
		CompoundLeft:     nil,
		CompoundRight:    nil,
		Comparison: &zbprotocol.TableSubqueryComparison{
			Op:       zbprotocol.QueryOperator_GREATER_THAN,
			Field:    field,
			Value:    valuStr,
			Ordering: qOrder,
		},
	}
}

func (q *QueryTextSearch) ToSubQuery(a, b string) *zbprotocol.TableSubQuery {
	valu := q.CompValue
	field := q.Field
	valuStr, qOrder := queryObjectTypify(valu)
	qOrder = zbprotocol.QueryOrdering_FULL_TEXT
	return &zbprotocol.TableSubQuery{
		IsCompound:       false,
		CompoundOperator: 0,
		CompoundLeft:     nil,
		CompoundRight:    nil,
		Comparison: &zbprotocol.TableSubqueryComparison{
			Op:       zbprotocol.QueryOperator_TEXT_SEARCH,
			Field:    field,
			Value:    valuStr,
			Ordering: qOrder,
		},
	}
}

func (q *QueryNotEquals) ToSubQuery(a, b string) *zbprotocol.TableSubQuery {
	valu := q.CompValue
	field := q.Field
	valuStr, qOrder := queryObjectTypify(valu)
	return &zbprotocol.TableSubQuery{
		IsCompound:       false,
		CompoundOperator: 0,
		CompoundLeft:     nil,
		CompoundRight:    nil,
		Comparison: &zbprotocol.TableSubqueryComparison{
			Op:       zbprotocol.QueryOperator_NOT_EQUALS,
			Field:    field,
			Value:    valuStr,
			Ordering: qOrder,
		},
	}
}

func (q *QueryEquals) ToSubQuery(a, b string) *zbprotocol.TableSubQuery {
	valu := q.CompValue
	field := q.Field
	valuStr, qOrder := queryObjectTypify(valu)
	return &zbprotocol.TableSubQuery{
		IsCompound:       false,
		CompoundOperator: 0,
		CompoundLeft:     nil,
		CompoundRight:    nil,
		Comparison: &zbprotocol.TableSubqueryComparison{
			Op:       zbprotocol.QueryOperator_EQUALS,
			Field:    field,
			Value:    valuStr,
			Ordering: qOrder,
		},
	}
}

func queryObjectTypify(valu interface{}) (string, zbprotocol.QueryOrdering) {
	var valuStr string
	qOrder := zbprotocol.QueryOrdering_LEXICOGRAPHIC
	switch valu.(type) {
	case int:
		qOrder = zbprotocol.QueryOrdering_REAL_NUMBERS
		valuStr = fmt.Sprintf("%d", valu.(int))
	case int64:
		qOrder = zbprotocol.QueryOrdering_REAL_NUMBERS
		valuStr = fmt.Sprintf("%d", valu.(int64))
	case int32:
		qOrder = zbprotocol.QueryOrdering_REAL_NUMBERS
		valuStr = fmt.Sprintf("%d", valu.(int32))
	case uint:
		qOrder = zbprotocol.QueryOrdering_REAL_NUMBERS
		valuStr = fmt.Sprintf("%d", valu.(uint))
	case uint8:
		qOrder = zbprotocol.QueryOrdering_REAL_NUMBERS
		valuStr = fmt.Sprintf("%d", valu.(uint8))
	case uint64:
		qOrder = zbprotocol.QueryOrdering_REAL_NUMBERS
		valuStr = fmt.Sprintf("%d", valu.(uint64))
	case uint32:
		qOrder = zbprotocol.QueryOrdering_REAL_NUMBERS
		valuStr = fmt.Sprintf("%d", valu.(uint32))
	case float64:
		qOrder = zbprotocol.QueryOrdering_REAL_NUMBERS
		valuStr = fmt.Sprintf("%f", valu.(float64))
	case float32:
		qOrder = zbprotocol.QueryOrdering_REAL_NUMBERS
		valuStr = fmt.Sprintf("%f", valu.(float32))
	case string:
		valuStr = valu.(string)
	default:
		valuStr = fmt.Sprintf("%v", valu.(string))
	}
	return valuStr, qOrder
}

func QNEq(field string, valu interface{}) *QueryNotEquals {
	return &QueryNotEquals{
		Field:     field,
		CompValue: valu,
	}
}

func QEq(field string, valu interface{}) *QueryEquals {
	return &QueryEquals{
		Field:     field,
		CompValue: valu,
	}
}

func QLte(field string, valu interface{}) *QueryLessThanEqual {
	return &QueryLessThanEqual{
		Field:     field,
		CompValue: valu,
	}
}

func QLt(field string, valu interface{}) *QueryLessThan {
	return &QueryLessThan{
		Field:     field,
		CompValue: valu,
	}
}

func QGt(field string, valu interface{}) *QueryGreaterThan {
	return &QueryGreaterThan{
		Field:     field,
		CompValue: valu,
	}
}

func QGte(field string, valu interface{}) *QueryGreaterThanEqual {
	return &QueryGreaterThanEqual{
		Field:     field,
		CompValue: valu,
	}
}

func QText(field string, queryStr string) *QueryTextSearch {
	return &QueryTextSearch{
		Field:     field,
		CompValue: queryStr,
	}
}
