package mysql

import (
	"database/sql"
	"encoding/json"
	"errors"
	"reflect"
	"strings"

	"github.com/uccu/go-stringify"
)

type Orm struct {
	db        *DB
	table     string
	rawQuery  bool
	query     string
	args      []interface{}
	dest      interface{}
	fields    []string
	err       error
	rawFields bool
	wk        []string
	wv        []interface{}
	sk        []string
	sv        []interface{}
}

func (v *Orm) Query(query string, args ...interface{}) *Orm {
	v.query = query
	v.args = args
	return v
}

func (v *Orm) RawQuery(b ...bool) *Orm {
	if len(b) > 0 {
		v.rawQuery = b[0]
	} else {
		v.rawQuery = true
	}
	return v
}

func (v *Orm) Dest(dest interface{}) *Orm {
	v.dest = dest
	return v
}

func (v *Orm) Fields(fields []string) *Orm {
	v.fields = fields
	return v
}

func (v *Orm) Field(fields ...string) *Orm {
	v.fields = fields
	return v
}

func (v *Orm) transfOrmFields() string {
	if len(v.fields) == 0 {
		return "*"
	}
	if v.rawFields {
		return strings.Join(v.fields, ",")
	}
	return "`" + strings.Join(v.fields, "`,`") + "`"
}

func (v *Orm) RawFields(r bool) *Orm {
	v.rawFields = r
	return v
}

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

func (v *Orm) Err() error {
	return v.err
}

func (v *Orm) Count() int64 {
	return v.RawFields(true).GetFieldInt("COUNT(1)")
}

func (v *Orm) setErr(e error) *Orm {
	v.err = e
	return v
}

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

func (v *Orm) WhereStru(s interface{}) *Orm {
	m, _ := json.Marshal(s)
	p := map[string]interface{}{}
	err := json.Unmarshal(m, &p)
	v.setErr(err)
	return v.Where(p)
}

func (v *Orm) Where(data map[string]interface{}) *Orm {

	if v.wk == nil {
		v.wk = make([]string, 0)
	}

	if v.wv == nil {
		v.wv = make([]interface{}, 0)
	}

	for k, vs := range data {
		v.wk = append(v.wk, k)
		v.wv = append(v.wv, vs)
	}

	return v
}

func (v *Orm) SetStru(s interface{}) *Orm {
	m, _ := json.Marshal(s)
	p := map[string]interface{}{}
	err := json.Unmarshal(m, &p)
	v.setErr(err)
	return v.Set(p)
}

func (v *Orm) Set(data map[string]interface{}) *Orm {

	if v.sk == nil {
		v.sk = make([]string, 0)
	}

	if v.sv == nil {
		v.sv = make([]interface{}, 0)
	}

	for k, vs := range data {
		v.sk = append(v.sk, k)
		v.sv = append(v.sv, vs)
	}

	return v
}

func (v *Orm) transfOrmQuery() string {

	where := ""
	if len(v.wk) > 0 {
		v.args = append(v.wv, v.args)
		where = " WHERE `" + strings.Join(v.wk, "`=? AND `") + "`=?"
	}

	set := ""
	if len(v.sk) > 0 {
		v.args = append(v.sv, v.args)
		set = " SET `" + strings.Join(v.sk, "`=?,`") + "`=?"
	}

	query := ""
	if v.query != "" {
		query = " " + v.query
	}

	return set + where + query
}

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

func (v *Orm) Update() (int64, error) {
	result, err := v.db.Exec("UPDATE "+v.table+v.transfOrmQuery(), v.args...)
	if err != nil {
		return 0, err
	}
	return result.RowsAffected()
}

func (v *Orm) Insert() (int64, error) {
	result, err := v.db.Exec("INSERT INTO "+v.table+v.transfOrmQuery(), v.args...)
	if err != nil {
		v.setErr(err)
		return 0, err
	}
	return result.LastInsertId()
}

func (v *Orm) Delete() (int64, error) {
	result, err := v.db.Exec("DELETE FROM "+v.table+v.transfOrmQuery(), v.args...)
	if err != nil {
		v.setErr(err)
		return 0, err
	}
	return result.RowsAffected()
}

var (
	NOT_PTR           = errors.New("not pass a pointer")
	NOT_SLI           = errors.New("not pass a slice")
	NOT_STRU          = errors.New("not pass a struct")
	NOT_STRU_IN_SLICE = errors.New("not pass a struct in slice")
	NIL_PTR           = errors.New("pass a nil pointer")
	NO_ROWS           = errors.New("no rows")
)

type column struct {
	Dest interface{}
}

func scanSlice(dest interface{}, rows *sql.Rows) error {

	value, err := getSlice(dest)
	if err != nil {
		return err
	}

	base, isPtr := getSliceBase(value)

	if base.Kind() != reflect.Struct {
		return NOT_STRU_IN_SLICE
	}

	columns, err := rows.Columns()
	if err != nil {
		return err
	}

	for rows.Next() {
		r := reflect.New(base)
		rv := stringify.GetReflectValue(r)
		rows.Scan(generateScanData(rv, columns)...)

		if isPtr {
			value.Set(reflect.Append(value, r))
		} else {
			value.Set(reflect.Append(value, rv))
		}
	}

	return nil

}

func scanOne(dest interface{}, rows *sql.Rows) error {

	value := reflect.ValueOf(dest)

	if value.Kind() != reflect.Ptr {
		return NOT_PTR
	}

	value = stringify.GetReflectValue(value)

	if value.Kind() != reflect.Struct {
		return NOT_STRU
	}

	columns, err := rows.Columns()
	if err != nil {
		return err
	}

	if !rows.Next() {
		return NO_ROWS
	}
	rows.Scan(generateScanData(value, columns)...)
	return nil
}

func generateScanData(rv reflect.Value, columns []string) []interface{} {

	s := []interface{}{}

	columnMap := map[string]*column{}
	for _, v := range columns {
		columnMap[v] = &column{}
	}

	for k := 0; k < rv.NumField(); k++ {
		name := rv.Type().Field(k).Tag.Get("db")
		if name == "" {
			continue
		}

		if column, ok := columnMap[name]; ok {
			f := rv.Field(k)
			if f.CanAddr() && f.CanInterface() {
				column.Dest = f.Addr().Interface()
			}
		}
	}

	for _, v := range columns {
		if columnMap[v].Dest == nil {
			s = append(s, nil)
		} else {
			s = append(s, columnMap[v].Dest)
		}
	}

	return s
}

func getSlice(dest interface{}) (reflect.Value, error) {
	value := reflect.ValueOf(dest)

	if value.Kind() != reflect.Ptr {
		return reflect.Value{}, NOT_PTR
	}

	value = stringify.GetReflectValue(value)

	if value.IsNil() {
		return reflect.Value{}, NIL_PTR
	}

	if value.Kind() != reflect.Slice {
		return reflect.Value{}, NOT_SLI
	}

	return value, nil
}

func getSliceBase(value reflect.Value) (base reflect.Type, isPtr bool) {

	slice := value.Type()
	base = slice.Elem()
	isPtr = base.Kind() == reflect.Ptr
	if isPtr {
		base = base.Elem()
	}
	return
}
