package tls

import (
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

// MySQL is tls core
type MySQL struct {
	*sqlx.DB
	Tag string // tag used for kernel operations to parse the field information of the structure based on this label
}

// Open is the same as sql.Open, but returns an *tls.MySQL instead.
func Open(dataSourceName string, tag string) (*MySQL, error) {
	db, err := sqlx.Open("mysql", dataSourceName)
	if err != nil {
		return nil, err
	}
	return &MySQL{DB: db, Tag: tag}, nil
}

// MustOpen is the same as sql.Open, but returns an *tls.MySQL instead and panics on error.
func MustOpen(dataSourceName string, tag string) *MySQL {
	db := sqlx.MustOpen("mysql", dataSourceName)
	return &MySQL{DB: db, Tag: tag}
}

// Engine generate this obj sql engine
func (mysql *MySQL) Engine(obj TableInterface) *Engine[TableInterface] {
	return &Engine[TableInterface]{
		MySQL:     mysql,
		obj:       obj,
		columns:   ParserColumns(obj, mysql.Tag, true),
		tableName: obj.TableName(),
	}
}
