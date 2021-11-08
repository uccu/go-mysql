package mysql

import (
	"reflect"
	"time"

	"github.com/uccu/go-mysql/field"
	"github.com/uccu/go-stringify"
)

func (v *Orm) startQuery(sql string) {
	v.StartQueryTime = time.Now()
	v.Sql = sql
}

func (v *Orm) afterQuery() {
	if v.db.afterQueryHandler != nil {
		v.db.afterQueryHandler(v)
	}
}

// 获取多条数据
func (v *Orm) Select() error {

	if len(v.table) > 0 {
		v.startQuery(v.transformSelectSql())
	} else {
		v.startQuery(v.transformQuery())
	}

	if len(v.orms) > 0 {
		for _, o := range v.orms {
			o.Select()
			v.Sql += " UNION "
			if o.unionAll {
				v.Sql += "ALL "
			}
			o.Select()
			v.Sql += "(" + o.Sql + ")"
		}
	}

	if v.b {
		return nil
	}

	rows, err := v.db.Query(v.Sql, v.GetArgs()...)
	if err != nil {
		v.setErr(err)
		return err
	}

	defer rows.Close()
	v.afterQuery()

	err = scanSlice(v.dest, rows)
	if err != nil {
		v.setErr(err)
		return err
	}

	return nil
}

// 获取单条数据
func (v *Orm) FetchOne() error {

	v.Limit(1)

	if len(v.table) > 0 {
		v.startQuery(v.transformSelectSql())
	} else {
		v.startQuery(v.transformQuery())
	}

	if v.b {
		return nil
	}

	rows, err := v.db.Query(v.Sql, v.GetArgs()...)
	if err != nil {
		v.setErr(err)
		return err
	}

	defer rows.Close()
	v.afterQuery()

	err = scanOne(v.dest, rows)
	if err != nil {
		v.setErr(err)
		return err
	}

	return nil
}

// 更新
func (v *Orm) Update() (int64, error) {

	if len(v.table) > 0 {
		v.startQuery(v.transformUpdateSql())
	} else {
		v.startQuery(v.transformQuery())
	}

	if v.b {
		return 0, nil
	}

	result, err := v.db.Exec(v.Sql, v.GetArgs()...)
	if err != nil {
		v.setErr(err)
		return 0, err
	}

	v.afterQuery()

	return result.RowsAffected()
}

// 插入
func (v *Orm) Insert() (int64, error) {

	if len(v.table) > 0 {
		v.startQuery(v.transformInsertSql())
	} else {
		v.startQuery(v.transformQuery())
	}

	if v.b {
		return 0, nil
	}

	result, err := v.db.Exec(v.Sql, v.GetArgs()...)
	if err != nil {
		v.setErr(err)
		return 0, err
	}
	v.afterQuery()

	return result.LastInsertId()
}

func (v *Orm) Replace() (int64, error) {

	if len(v.table) > 0 {
		v.startQuery(v.transformReplaceSql())
	} else {
		v.startQuery(v.transformQuery())
	}

	if v.b {
		return 0, nil
	}

	result, err := v.db.Exec(v.Sql, v.GetArgs()...)
	if err != nil {
		v.setErr(err)
		return 0, err
	}
	v.afterQuery()

	return result.LastInsertId()
}

// 删除
func (v *Orm) Delete() (int64, error) {

	if len(v.table) > 0 {
		v.startQuery(v.transformDeleteSql())
	} else {
		v.startQuery(v.transformQuery())
	}

	if v.b {
		return 0, nil
	}

	result, err := v.db.Exec(v.Sql, v.GetArgs()...)
	if err != nil {
		v.setErr(err)
		return 0, err
	}
	v.afterQuery()
	return result.RowsAffected()
}

// 获取单个字段的值
func (v *Orm) GetField(name interface{}) error {

	v.Field(name)
	v.Limit(1)

	if len(v.table) > 0 {
		v.startQuery(v.transformSelectSql())
	} else {
		v.startQuery(v.transformQuery())
	}

	if v.b {
		return nil
	}

	row := v.db.QueryRow(v.Sql, v.GetArgs()...)
	if row.Err() != nil {
		v.setErr(row.Err())
		return row.Err()
	}

	v.afterQuery()

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

	if len(v.table) > 0 {
		v.startQuery(v.transformSelectSql())
	} else {
		v.startQuery(v.transformQuery())
	}

	if v.b {
		return nil
	}

	rows, err := v.db.Query(v.Sql, v.GetArgs()...)
	if err != nil {
		v.setErr(err)
		return err
	}
	defer rows.Close()
	v.afterQuery()

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

func (v *Orm) Count(f ...string) int64 {

	if len(f) > 0 {
		k := transformToKey(f[0])
		return v.GetFieldInt(field.NewMutiField("COUNT(%t)", field.NewField(k.Name).SetAlias(k.Alias).SetTable(k.Parent)))
	}
	return v.GetFieldInt(field.NewRawField("COUNT(1)"))
}

func (v *Orm) Exist() bool {
	return v.GetFieldInt(field.NewRawField("1")) != 0
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
