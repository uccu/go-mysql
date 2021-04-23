package mysql

import (
	"reflect"

	"github.com/uccu/go-stringify"
)

// 获取多条数据
func (v *Orm) Select() error {

	sql := v.query
	if !v.rawQuery {
		sql = "SELECT " + v.transfOrmFields() + " FROM " + v.table + v.transfOrmQuery()
	}

	rows, err := v.db.Query(sql, v.args...)
	defer rows.Close()
	if err != nil {
		v.setErr(err)
		return err
	}

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
	if !v.rawQuery {
		sql = "SELECT " + v.transfOrmFields() + " FROM " + v.table + v.transfOrmQuery() + " LIMIT 1"
	}
	rows, err := v.db.Query(sql, v.args...)
	defer rows.Close()
	if err != nil {
		v.setErr(err)
		return err
	}

	err = scanOne(v.dest, rows)
	if err != nil {
		v.setErr(err)
		return err
	}

	return nil
}

// 更新
func (v *Orm) Update() (int64, error) {
	result, err := v.db.Exec("UPDATE "+v.table+v.transfOrmQuery(), v.args...)
	if err != nil {
		return 0, err
	}
	return result.RowsAffected()
}

// 插入
func (v *Orm) Insert() (int64, error) {
	result, err := v.db.Exec("INSERT INTO "+v.table+v.transfOrmQuery(), v.args...)
	if err != nil {
		v.setErr(err)
		return 0, err
	}
	return result.LastInsertId()
}

// 删除
func (v *Orm) Delete() (int64, error) {
	result, err := v.db.Exec("DELETE FROM "+v.table+v.transfOrmQuery(), v.args...)
	if err != nil {
		v.setErr(err)
		return 0, err
	}
	return result.RowsAffected()
}

// 获取单个字段的值
func (v *Orm) GetField(name string) error {
	v.Field(name)
	sql := "SELECT " + v.transfOrmFields() + " FROM " + v.table + v.transfOrmQuery() + " LIMIT 1"
	err := v.db.QueryRow(sql, v.args...).Scan(v.dest)

	if err != nil {
		v.setErr(err)
		return err
	}
	return nil
}

// 获取单个字段的值的slice
func (v *Orm) GetFields(name string) error {
	v.Field(name)
	sql := "SELECT " + v.transfOrmFields() + " FROM " + v.table + v.transfOrmQuery()
	rows, err := v.db.Query(sql, v.args...)
	defer rows.Close()

	if err != nil {
		v.setErr(err)
		return err
	}

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

func (v *Orm) GetFieldString(name string) string {
	var data string
	v.Dest(&data).GetField(name)
	return data
}

func (v *Orm) GetFieldInt(name string) int64 {
	var data int64
	v.Dest(&data).GetField(name)
	return data
}

func (v *Orm) Count(val ...string) int64 {
	field := "COUNT(1)"
	if len(val) > 0 {
		field = val[0]
	}
	return v.RawFields(true).GetFieldInt(field)
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
