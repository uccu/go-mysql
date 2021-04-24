package mysql

import "database/sql"

type DB struct {
	*sql.DB
	prefix string
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

func (v *DB) GetPrefix() string {
	return v.prefix
}

func Open(driverName, dataSourceName string) (*DB, error) {
	d, err := sql.Open(driverName, dataSourceName)
	if err != nil {
		return nil, err
	}
	return &DB{DB: d}, nil
}
