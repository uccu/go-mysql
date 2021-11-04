package mysql

import "database/sql"

type DB struct {
	*sql.DB
	prefix            string
	subTable          func(table string, n int64) string
	errHandler        func(error)
	afterQueryHandler func(*Orm)
}

func (v *DB) GetOrm(name ...string) *Orm {
	var table string
	if len(name) > 0 {
		table = name[0]
	}
	return &Orm{table: table, db: v}
}

func (v *DB) Prefix(p string) *DB {
	v.prefix = p
	return v
}

func (v *DB) ErrHandle(p func(error)) *DB {
	v.errHandler = p
	return v
}

func (v *DB) AfterQueryHandle(p func(*Orm)) *DB {
	v.afterQueryHandler = p
	return v
}

func (v *DB) SubTable(p func(table string, n int64) string) *DB {
	v.subTable = p
	return v
}

func Open(driverName, dataSourceName string) (*DB, error) {
	d, err := sql.Open(driverName, dataSourceName)
	if err != nil {
		return nil, err
	}
	return &DB{DB: d}, nil
}
