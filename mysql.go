package mysql

import (
	"database/sql"

	_ "github.com/go-sql-driver/mysql"
	"github.com/uccu/go-mysql/mx"
	"github.com/uccu/go-mysql/table"
)

type DB struct {
	*sql.DB
	prefix            string
	suffix            func(interface{}) string
	errHandler        func(error)
	afterQueryHandler func(*Orm)
}

func (v *DB) GetOrm(s ...interface{}) *Orm {
	tables := mx.Tables{}
	for _, t := range s {
		if t, ok := t.(mx.Table); ok {
			tables = append(tables, t)
			break
		}
		if t, ok := t.(string); ok {
			tables = append(tables, v.NewTable(t))
			break
		}
	}

	if len(tables) > 0 {
		return &Orm{table: tables, mix: make(mx.Mixs, 0), db: v}
	} else {
		return &Orm{table: make(mx.Tables, 0), mix: make(mx.Mixs, 0), db: v}
	}

}

func (v *DB) WithPrefix(p string) *DB {
	v.prefix = p
	return v
}

func (v *DB) WithSuffix(p func(interface{}) string) *DB {
	v.suffix = p
	return v
}

func (v *DB) WithErrHandler(p func(error)) *DB {
	v.errHandler = p
	return v
}

func (v *DB) WithAfterQueryHandler(p func(*Orm)) *DB {
	v.afterQueryHandler = p
	return v
}

func (v *DB) NewTable(name string, as ...string) *table.Table {
	t := table.NewTable(name, v.prefix, v.suffix)
	if len(as) > 0 {
		t.As = as[0]
	}
	return t
}

func Open(driverName, dataSourceName string) (*DB, error) {
	d, err := sql.Open(driverName, dataSourceName)
	if err != nil {
		return nil, err
	}
	return &DB{DB: d}, nil
}
