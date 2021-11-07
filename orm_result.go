package mysql

import (
	"reflect"
	"regexp"
	"time"

	"github.com/uccu/go-stringify"
)

// 获取多条数据
func (v *Orm) Select() error {

	sql := v.query
	if v.tables != nil {
		sql = v.transformSelectSql()
	}
	v.StartQueryTime = time.Now()
	v.Sql = sql
	rows, err := v.db.Query(sql, v.args...)
	if err != nil {
		v.setErr(err)
		return err
	}

	if v.db.afterQueryHandler != nil {
		v.db.afterQueryHandler(v)
	}

	defer rows.Close()

	err = scanSlice(v.dest, rows)
	if err != nil {
		v.setErr(err)
		return err
	}

	return nil
}

// 获取单条数据
func (v *Orm) FetchOne() error {
	sql := v.query
	if v.tables != nil {
		sql = v.transformSelectSql()
		r, _ := regexp.Compile(`(?i)(LIMIT +\d+ *)$`)
		if r.MatchString(sql) {
			sql = r.ReplaceAllString(sql, "LIMIT 1")
		} else {
			sql += " LIMIT 1"
		}
	}
	v.StartQueryTime = time.Now()
	v.Sql = sql
	rows, err := v.db.Query(sql, v.args...)
	if err != nil {
		v.setErr(err)
		return err
	}
	if v.db.afterQueryHandler != nil {
		v.db.afterQueryHandler(v)
	}
	defer rows.Close()

	err = scanOne(v.dest, rows)
	if err != nil {
		v.setErr(err)
		return err
	}

	return nil
}

// 更新
func (v *Orm) Update() (int64, error) {
	sql := "UPDATE " + v.tables.GetQuery() + v.transformQuery()
	result, err := v.db.Exec(sql, v.args...)
	if err != nil {
		v.setErr(err)
		return 0, err
	}
	return result.RowsAffected()
}

// 插入
func (v *Orm) Insert() (int64, error) {
	sql := "INSERT INTO " + v.tables.GetQuery() + v.transformQuery()
	result, err := v.db.Exec(sql, v.args...)
	if err != nil {
		v.setErr(err)
		return 0, err
	}
	return result.LastInsertId()
}

// 删除
func (v *Orm) Delete() (int64, error) {
	sql := "DELETE FROM " + v.tables.GetQuery() + v.transformQuery()
	result, err := v.db.Exec(sql, v.args...)
	if err != nil {
		v.setErr(err)
		return 0, err
	}
	return result.RowsAffected()
}

// 获取单个字段的值
func (v *Orm) GetField(name interface{}) error {
	v.Field(name)
	sql := v.transformSelectSql()
	r, _ := regexp.Compile(`(?i)(LIMIT +\d+ *)$`)
	if r.MatchString(sql) {
		sql = r.ReplaceAllString(sql, "LIMIT 1")
	} else {
		sql += " LIMIT 1"
	}
	v.StartQueryTime = time.Now()
	v.Sql = sql
	row := v.db.QueryRow(sql, v.args...)
	if row.Err() != nil {
		v.setErr(row.Err())
		return row.Err()
	}

	if v.db.afterQueryHandler != nil {
		v.db.afterQueryHandler(v)
	}

	err := row.Scan(v.dest)
	if err != nil {
		v.setErr(err)
		return err
	}

	return nil
}

// 获取单个字段的值的slice
func (v *Orm) GetFields(name string) error {
	v.Field(name)
	sql := v.transformSelectSql()
	v.StartQueryTime = time.Now()
	v.Sql = sql
	rows, err := v.db.Query(sql, v.args...)
	if err != nil {
		v.setErr(err)
		return err
	}
	if v.db.afterQueryHandler != nil {
		v.db.afterQueryHandler(v)
	}
	defer rows.Close()

	value, err := getSlice(v.dest)
	if err != nil {
		v.setErr(err)
		return err
	}
	base, isPtr := getSliceBase(value)

	for rows.Next() {
		b := reflect.New(base)
		bv := stringify.GetReflectValue(b)
		if bv.CanAddr() && bv.CanInterface() {
			rows.Scan(bv.Addr().Interface())
		}
		if isPtr {
			value.Set(reflect.Append(value, b))
		} else {
			value.Set(reflect.Append(value, bv))
		}
	}

	return nil
}

func (v *Orm) GetFieldString(f interface{}) string {
	var data string
	v.Dest(&data).GetField(f)
	return data
}

func (v *Orm) GetFieldInt(f interface{}) int64 {
	var data int64
	v.Dest(&data).GetField(f)
	return data
}

func (v *Orm) Count(f ...interface{}) int64 {
	if len(f) > 0 {
		k := v.transformField(f[0])
		if k != nil {
			return v.GetFieldInt(MutiField("COUNT(?)", k))
		}
	}
	return v.GetFieldInt(Raw("COUNT(1)"))
}

func (v *Orm) Exist() bool {
	return v.GetFieldInt(Raw("1")) != 0
}

func (v *Orm) GetFieldsString(name string) []string {
	data := []string{}
	v.Dest(&data).GetFields(name)
	return data
}

func (v *Orm) GetFieldsInt(name string) []int64 {
	data := []int64{}
	v.Dest(&data).GetFields(name)
	return data
}
