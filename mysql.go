package mysql

import "database/sql"

const (
	Nil = iota
)

type DB struct {
	*sql.DB
}

func (v *DB) Table(name ...string) *Orm {
	var table string
	if len(name) > 0 {
		table = name[0]
	}
	return &Orm{table: table, db: v}
}

func Open(driverName, dataSourceName string) (*DB, error) {
	d, err := sql.Open(driverName, dataSourceName)
	if err != nil {
		return nil, err
	}
	return &DB{d}, nil
}
