package mysql

import (
	"database/sql"

	_ "github.com/go-sql-driver/mysql"
)

type DB struct {
	*sql.DB
	prefix            string
	suffix            func(interface{}) string
	errHandler        func(error)
	afterQueryHandler func(*Orm)
}

func (db *DB) Table(s ...interface{}) *Orm {
	return db.Default().Table(s...)
}

func (db *DB) Default() *Orm {
	return &Orm{db: db}
}

func (db *DB) WithPrefix(p string) *DB {
	db.prefix = p
	return db
}

func (db *DB) WithSuffix(p func(interface{}) string) *DB {
	db.suffix = p
	return db
}

func (db *DB) WithErrHandler(p func(error)) *DB {
	db.errHandler = p
	return db
}

func (db *DB) WithAfterQueryHandler(p func(*Orm)) *DB {
	db.afterQueryHandler = p
	return db
}

func Open(driverName, dataSourceName string) (*DB, error) {
	d, err := sql.Open(driverName, dataSourceName)
	if err != nil {
		return nil, err
	}
	return &DB{DB: d}, nil
}
