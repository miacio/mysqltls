package tls

import (
	"database/sql"
	"fmt"
	"strings"
)

// TableInterface
type TableInterface interface {
	TableName() string    // TableName returnes the table name
	PrimaryKey() []string // PrimaryKey returnes the table id columns
}

// Engine sql builder for generating single table operations
type Engine[T TableInterface] struct {
	*MySQL
	obj       T        // core obj build current table sql
	columns   []string // the obj all columns
	tableName string   // the obj table name
}

// FindById based on PrimaryKey()
// the query statement based on PrimaryKey will be automatically renewed
func (e *Engine[T]) FindById(dest T, params ...any) error {
	sqlTmpl := "SELECT %s FROM %s WHERE %s"
	where := e.primaryKeyWhere()
	sql := fmt.Sprintf(sqlTmpl, strings.Join(e.columns, ","), e.tableName, where)
	return e.Get(dest, sql, params...)
}

// DelById based on PrimaryKey()
// the deletion statement based on PrimaryKey will be automatically renewed
func (e *Engine[T]) DeleteById(params ...any) (sql.Result, error) {
	sqlTmpl := "DELETE FROM %s WHERE %s"
	where := e.primaryKeyWhere()
	sql := fmt.Sprintf(sqlTmpl, e.tableName, where)
	return e.Exec(sql, params...)
}

// UpdateById based on PrimaryKey()
// the modification statement based on PrimaryKey will be automatically renewed
func (e *Engine[T]) UpdateById(dest T) (sql.Result, error) {
	sqlTmpl := "UPDATE %s SET %s WHERE %s"
	where := e.primaryKeyWhere()

	paramMap, err := ParserTagToMap(dest, e.Tag)
	if err != nil {
		return nil, err
	}
	setColumns, vals := ParserClause(paramMap, true, e.obj.PrimaryKey()...)
	if setColumns == nil {
		return nil, nil
	}
	for i := range setColumns {
		setColumns[i] = setColumns[i] + " = ?"
	}
	setClause := strings.Join(setColumns, ",")
	sql := fmt.Sprintf(sqlTmpl, e.tableName, setClause, where)

	for _, primaryKey := range e.obj.PrimaryKey() {
		vals = append(vals, paramMap[primaryKey])
	}
	return e.Exec(sql, vals...)
}

// BatchInsert
func (e *Engine[T]) BatchInsert(dests ...T) (sql.Result, error) {
	sqlTmpl := "INSERT INTO %s (%s) VALUES (%s)"

	namedColumns := ParserColumns(e.obj, e.Tag, false)
	valColumns := make([]string, 0)
	for i := range namedColumns {
		valColumns = append(valColumns, ":"+namedColumns[i])
	}
	sql := fmt.Sprintf(sqlTmpl, e.tableName, strings.Join(e.columns, ","), strings.Join(valColumns, ","))
	params := make([]map[string]interface{}, 0)
	for i := range dests {
		param, err := ParserTagToMap(dests[i], e.Tag)
		if err != nil {
			return nil, err
		}
		params = append(params, param)
	}
	return e.NamedExec(sql, params)
}

// primaryKeyWhere clause generation
func (e *Engine[T]) primaryKeyWhere() string {
	ids := e.obj.PrimaryKey()
	if ids == nil {
		panic("table primary key columns is empty")
	}
	whereColumns := make([]string, 0)
	for _, idColumn := range ids {
		whereColumns = append(whereColumns, fmt.Sprintf("%s = ?", KeywordTo(idColumn)))
	}
	where := strings.Join(whereColumns, " AND ")
	return where
}
