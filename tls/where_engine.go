package tls

import (
	"database/sql"
	"fmt"
	"strings"

	"github.com/jmoiron/sqlx"
)

type WhereEngine[T TableInterface] struct {
	*Engine[T]

	whereClause []Clause // where clause
}

// Where
func (e *Engine[T]) Where() *WhereEngine[T] {
	selectEngine := WhereEngine[T]{
		e,
		make([]Clause, 0),
	}
	return &selectEngine
}

// And
func (se *WhereEngine[T]) And(condition string, params ...interface{}) *WhereEngine[T] {
	return se.whereAppend(false, condition, params...)
}

// AndIn
func (se *WhereEngine[T]) AndIn(condition string, params ...interface{}) *WhereEngine[T] {
	sql, params, err := sqlx.In(condition, params...)
	if err != nil {
		panic(err)
	}
	return se.whereAppend(false, sql, params...)
}

// OrIn
func (se *WhereEngine[T]) OrIn(condition string, params ...interface{}) *WhereEngine[T] {
	sql, params, err := sqlx.In(condition, params...)
	if err != nil {
		panic(err)
	}
	return se.newAppend().whereAppend(false, sql, params...)
}

// Or
func (se *WhereEngine[T]) Or(condition string, params ...interface{}) *WhereEngine[T] {
	return se.newAppend().whereAppend(false, condition, params...)
}

func (se *WhereEngine[T]) newAppend() *WhereEngine[T] {
	if len(se.whereClause) > 0 {
		clause := se.whereClause[len(se.whereClause)-1]
		clause.End = true
		se.whereClause[len(se.whereClause)-1] = clause
	}
	return se
}

// where clause append
func (se *WhereEngine[T]) whereAppend(end bool, condition string, params ...interface{}) *WhereEngine[T] {
	if strings.Trim(condition, " ") == "" {
		return se
	}
	var clause Clause
	if se.whereClause == nil || len(se.whereClause) == 0 {
		clause = NewClause()
		clause.Condition = append(clause.Condition, condition)
		clause.Params = append(clause.Params, params...)
		clause.End = end
		se.whereClause = append(se.whereClause, clause)
		return se
	}
	clause = se.whereClause[len(se.whereClause)-1]
	if !clause.End {
		clause.Condition = append(clause.Condition, condition)
		clause.Params = append(clause.Params, params...)
		se.whereClause[len(se.whereClause)-1] = clause
	} else {
		newClause := NewClause()
		newClause.Condition = append(newClause.Condition, condition)
		newClause.Params = append(newClause.Params, params...)
		newClause.End = end
		se.whereClause = append(se.whereClause, newClause)
	}
	return se
}

// whereRead
func (se *WhereEngine[T]) whereRead() (string, []interface{}) {
	sqlChain := make([]string, 0)
	params := make([]interface{}, 0)
	for i := range se.whereClause {
		sqlChain = append(sqlChain, strings.Join(se.whereClause[i].Condition, " and "))
		params = append(params, se.whereClause[i].Params...)
	}
	if len(sqlChain) > 1 {
		for i := range sqlChain {
			sqlChain[i] = "(" + sqlChain[i] + ")"
		}
	}
	return strings.Join(sqlChain, " or "), params
}

// Find
func (se *WhereEngine[T]) Find(dest any, columns ...string) error {
	sqlTmpl := "SELECT %s FROM %s"
	if columns != nil {
		for i := range columns {
			columns[i] = KeywordTo(columns[i])
		}
	} else {
		columns = se.columns
	}
	sql := fmt.Sprintf(sqlTmpl, strings.Join(columns, ","), se.tableName)
	whereCondition, params := se.whereRead()
	if strings.Trim(whereCondition, " ") != "" {
		sql = sql + " WHERE " + whereCondition
	}
	return se.DB.Select(dest, sql, params...)
}

// Get
func (se *WhereEngine[T]) Get(dest any, columns ...string) error {
	sqlTmpl := "SELECT %s FROM %s"
	if columns != nil {
		for i := range columns {
			columns[i] = KeywordTo(columns[i])
		}
	} else {
		columns = se.columns
	}
	sql := fmt.Sprintf(sqlTmpl, strings.Join(columns, ","), se.tableName)
	whereCondition, params := se.whereRead()
	if strings.Trim(whereCondition, " ") != "" {
		sql = sql + " WHERE " + whereCondition
	}
	return se.DB.Get(dest, sql, params...)
}

// Count
func (se *WhereEngine[T]) Count() (int, error) {
	sqlTmpl := "SELECT COUNT(1) FROM %s"
	sql := fmt.Sprintf(sqlTmpl, se.tableName)
	whereCondition, params := se.whereRead()
	if strings.Trim(whereCondition, " ") != "" {
		sql = sql + " WHERE " + whereCondition
	}
	var result int
	err := se.DB.Get(&result, sql, params...)
	return result, err
}

// Delete
func (se *WhereEngine[T]) Delete() (sql.Result, error) {
	sqlTmpl := "DELETE FROM %s"
	sql := fmt.Sprintf(sqlTmpl, se.tableName)
	whereCondition, params := se.whereRead()
	if strings.Trim(whereCondition, " ") != "" {
		sql = sql + " WHERE " + whereCondition
	}
	return se.DB.Exec(sql, params...)
}

// Update
func (se *WhereEngine[T]) Update(conditions string, params ...interface{}) (sql.Result, error) {
	sqlTmpl := "UPDATE %s SET %s"
	sql := fmt.Sprintf(sqlTmpl, se.tableName, conditions)
	whereCondition, whereParams := se.whereRead()
	if strings.Trim(whereCondition, " ") != "" {
		sql = sql + " WHERE " + whereCondition
	}
	updateParams := params
	updateParams = append(updateParams, whereParams...)
	fmt.Println(sql)
	fmt.Println(updateParams...)
	return se.Exec(sql, updateParams...)
}
