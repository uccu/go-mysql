package mysql

import (
	"database/sql"
	"errors"
	"reflect"

	"github.com/uccu/go-stringify"
)

var (
	NOT_PTR           = errors.New("not pass a pointer")
	NOT_SLI           = errors.New("not pass a slice")
	NOT_STRU          = errors.New("not pass a struct")
	NOT_STRU_IN_SLICE = errors.New("not pass a struct in slice")
	NIL_PTR           = errors.New("pass a nil pointer")
	ODD_PARAM         = errors.New("Odd number of parameters")
	NO_DB_TAG         = errors.New("no db tag")
	NO_ROWS           = sql.ErrNoRows
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

	loopStruct(rv, func(v reflect.Value, s reflect.StructField) bool {
		if name := s.Tag.Get("db"); name != "" {
			if column, ok := columnMap[name]; ok {
				if v.CanAddr() && v.CanInterface() {
					column.Dest = v.Addr().Interface()
				}
			}
			return true
		}
		return false
	})

	for _, v := range columns {
		s = append(s, columnMap[v].Dest)
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

func loopStructType(val reflect.Type, f func(s reflect.StructField) bool) {
	if val.Kind() != reflect.Struct {
		return
	}
	for k := 0; k < val.NumField(); k++ {
		ft := val.Field(k).Type
		for ft.Kind() == reflect.Ptr || ft.Kind() == reflect.Interface {
			ft = ft.Elem()
		}

		if !f(val.Field(k)) && ft.Kind() == reflect.Struct {
			loopStructType(ft, f)
		}

	}
}

func loopStruct(val reflect.Value, f func(v reflect.Value, s reflect.StructField) bool) {
	if val.Kind() != reflect.Struct {
		return
	}
	for k := 0; k < val.NumField(); k++ {
		ft := val.Field(k)
		for ft.Kind() == reflect.Ptr || ft.Kind() == reflect.Interface {
			if ft.Kind() == reflect.Ptr {
				if ft.Type().Elem().Kind() == reflect.Struct {
					if ft.Elem().Kind() == reflect.Invalid {
						ft.Set(reflect.New(ft.Type().Elem()))
					}
				}
			}
			ft = ft.Elem()
		}
		if !f(val.Field(k), val.Type().Field(k)) && ft.Kind() == reflect.Struct {
			loopStruct(ft, f)
		}
	}
}

func removeRep(s []string) []string {
	r := []string{}
	t := map[string]bool{}
	for _, e := range s {
		l := len(t)
		t[e] = false
		if len(t) != l {
			r = append(r, e)
		}
	}
	return r
}
