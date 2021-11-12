package mysql

import (
	"database/sql"
	"errors"
	"reflect"
	"regexp"
	"strings"

	"github.com/uccu/go-mysql/field"
	"github.com/uccu/go-mysql/mix"
	"github.com/uccu/go-mysql/mx"
	"github.com/uccu/go-mysql/table"
	"github.com/uccu/go-stringify"
)

var (
	ErrNotPointer        = errors.New("not pass a pointer")
	ErrNotSlice          = errors.New("not pass a slice")
	ErrNotStruct         = errors.New("not pass a struct")
	ErrNotStructInSlice  = errors.New("not pass a struct in slice")
	ErrNilPointer        = errors.New("pass a nil pointer")
	ErrOddNumberOfParams = errors.New("odd number of parameters")
	ErrNoContainer       = errors.New("no container")
	ErrType              = errors.New("error type")
	ErrNoRows            = sql.ErrNoRows
)

type sqlType byte
type resultType byte

const (
	SQL_SELECT sqlType = iota
	SQL_DELETE
	SQL_UPDATE
	SQL_INSERT
	SQL_REPLACE
)

const (
	RESULT_LAST_INSERT_ID resultType = iota
	RESULT_ROWS_AFFECTED
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
		return ErrNotStructInSlice
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
		return ErrNotPointer
	}

	value = stringify.GetReflectValue(value)

	if value.Kind() != reflect.Struct {
		return ErrNotStruct
	}

	columns, err := rows.Columns()
	if err != nil {
		return err
	}

	if !rows.Next() {
		return ErrNoRows
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
		return reflect.Value{}, ErrNotPointer
	}

	value = stringify.GetReflectValue(value)

	if value.IsNil() {
		return reflect.Value{}, ErrNilPointer
	}

	if value.Kind() != reflect.Slice {
		return reflect.Value{}, ErrNotSlice
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
		if f(ft, val.Type().Field(k)) {
			continue
		}

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
		if ft.Kind() == reflect.Struct {
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

type key struct {
	Alias  string
	Name   string
	Parent string
}

type keyList []*key

func transformToKeyList(i string) keyList {

	list := keyList{}
	slist := stringify.ToStringSlice(i, ",")

	for _, s := range slist {

		s = strings.Trim(s, " ")
		s = strings.ReplaceAll(s, "`", "")

		sli := stringify.ToStringSlice(s, ".")

		k := &key{}
		names := sli[0]
		if len(sli) > 1 {
			k.Parent = sli[0]
			names = sli[1]
		}

		sli = stringify.ToStringSlice(names, " ")
		k.Name = sli[0]
		if len(sli) > 1 {
			k.Alias = sli[len(sli)-1]
		}

		list = append(list, k)
	}

	return list
}

func transformToKey(i string) *key {
	return transformToKeyList(i)[0]
}

func transformStructToMixs(s interface{}, tagName string) mx.Mixs {

	p := mx.Mixs{}
	rv := stringify.GetReflectValue(s)

	loopStruct(rv, func(v reflect.Value, s reflect.StructField) bool {
		db := s.Tag.Get("db")
		dbset := s.Tag.Get(tagName)
		if dbset != "" {
			db = dbset
		}
		if db == "-" {
			return true
		}
		if db == "" {
			return false
		}
		if v.CanInterface() {
			p = append(p, Mix("%t=?", field.NewField(db).ToMix(), v.Interface()))

			return true
		}
		return false
	})

	if len(p) == 0 {
		return nil
	}
	return p
}

func transformMapToMixs(s map[string]interface{}) mx.Mixs {
	p := mx.Mixs{}
	for k, v := range s {
		p = append(p, Mix("%t=?", field.NewField(k).ToMix(), v))
	}
	return p
}

func transformSliceToMixs(s ...interface{}) mx.Mixs {
	p := mx.Mixs{}
	for k := 0; k < len(s); k += 2 {
		p = append(p, Mix("%t=?", field.NewField(s[k].(string)).ToMix(), s[k+1]))
	}
	return p
}

func transformToMixs(tagName string, s ...interface{}) (mx.Mixs, error) {
	var mixs mx.Mixs
	if len(s) == 1 {
		if s, ok := s[0].(mx.Mix); ok {
			return mx.Mixs{s}, nil
		}
		rv := stringify.GetReflectValue(s[0])
		if rv.Kind() == reflect.Struct {
			mixs = transformStructToMixs(s[0], tagName)
		} else if rv.Kind() == reflect.Map {
			mixs = transformMapToMixs(s[0].(map[string]interface{}))
		}
	} else {
		if len(s)%2 == 1 {
			return nil, ErrOddNumberOfParams
		}
		mixs = transformSliceToMixs(s...)
	}
	return mixs, nil
}

func Field(f string) *field.Field {
	k := transformToKey(f)
	return field.NewField(k.Name).SetAlias(k.Alias).SetTable(k.Parent)
}

func FieldMix(f string) mx.Mix {
	return Field(f).ToMix()
}

func Table(f string) *table.Table {
	k := transformToKey(f)
	return table.NewTable(k.Name).SetAlias(k.Alias).SetDBName(k.Parent)
}

func Raw(f string) *mix.Raw {
	return mix.NewRawMix(f)
}

func Mix(q string, f ...interface{}) *mix.Mix {
	r := regexp.MustCompile(`(?i)%t|\?`)
	loc := r.FindAllStringIndex(q, -1)

	k := 0
	mixs := mx.Mixs{}
	args := []interface{}{}

	for _, si := range loc {
		if q[si[0]:si[1]] == "%t" {
			if v, ok := f[k].(mx.Mix); ok {
				mixs = append(mixs, v)
				args = append(args, v.GetArgs()...)
			} else if v, ok := f[k].(mx.Field); ok {
				mixs = append(mixs, v.ToMix())
			} else if v, ok := f[k].(string); ok {
				mixs = append(mixs, Field(v).ToMix())
			} else {
				mixs = append(mixs, Raw("NULL"))
			}
		} else {
			args = append(args, f[k])
		}
		k++
	}
	return mix.NewMix(q, mixs, args)
}
