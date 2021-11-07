package mysql

import (
	"reflect"
	"time"

	"github.com/uccu/go-mysql/field"
	"github.com/uccu/go-mysql/mix"
	"github.com/uccu/go-mysql/mx"
	"github.com/uccu/go-stringify"
)

type Orm struct {
	db      *DB
	table   mx.Tables
	mix     mx.Mixs
	fields  mx.Fields
	mixType string

	dest interface{}
	err  error

	StartQueryTime time.Time
	Sql            string
}

func (v *Orm) Query(query string, args ...interface{}) *Orm {
	v.addMix(mix.NewMix(" "+query, args...))
	return v
}

func (v *Orm) addMix(m mx.Mix, typs ...string) *Orm {

	if v.mix == nil {
		v.mix = make(mx.Mixs, 0)
	}

	if len(typs) > 0 && v.mixType != typs[0] {
		typ := typs[0]
		if typ == "where" {
			v.mix = append(v.mix, mix.NewMix("WHERE"))
		} else if typ == "set" {
			v.mix = append(v.mix, mix.NewMix("SET"))
		} else if typ == "group" {
			v.mix = append(v.mix, mix.NewMix("GROUP BY"))
		} else if typ == "limit" {
			v.mix = append(v.mix, mix.NewMix("LIMIT"))
		} else if typ == "order" {
			v.mix = append(v.mix, mix.NewMix("ORDER BY"))
		} else if typ == "having" {
			v.mix = append(v.mix, mix.NewMix("HAVING"))
		}
		v.mixType = typ
	}

	v.mix = append(v.mix, m)
	return v
}

func (v *Orm) Dest(dest interface{}) *Orm {
	v.dest = dest
	return v
}

func (v *Orm) addField(field mx.Field) *Orm {
	if v.fields == nil {
		v.fields = make(mx.Fields, 0)
	}
	v.fields = append(v.fields, field)
	return v
}

func (v *Orm) transformDestToField() mx.Fields {
	val := stringify.GetReflectValue(v.dest).Type()
	if val.Kind() == reflect.Slice {
		val = val.Elem()
		for val.Kind() == reflect.Ptr {
			val = val.Elem()
		}
	}

	if val.Kind() == reflect.Struct {
		fields := []string{}
		loopStructType(val, func(s reflect.StructField) bool {
			name := s.Tag.Get("db")
			if name != "" {
				if name != "-" {
					fields = append(fields, name)
				}
				return true
			}
			return false
		})
		fields = removeRep(fields)
		if len(fields) == 0 {
			return mx.Fields{field.NewRawField("1")}
		}

		keys := mx.Fields{}
		for _, f := range fields {
			keys = append(keys, field.NewField(f))
		}
		return keys
	}

	return mx.Fields{field.NewRawField("1")}
}

func (v *Orm) transformFields() string {
	if len(v.fields) == 0 {
		if v.dest != nil {
			fields := v.transformDestToField()
			fields.With(mx.WithBackquote)
			if len(v.table) > 1 {
				fields.With(mx.WithTable)
			}
			return fields.GetQuery()
		}
		return "*"
	}

	if len(v.table) > 1 {
		v.fields.With(mx.WithTable)
	}
	v.fields.With(mx.WithBackquote)
	return v.fields.GetQuery()
}

func (v *Orm) transformSelectSql() string {
	return "SELECT " + v.transformFields() + " FROM " + v.transformTable() + " " + v.transformQuery()
}

func (v *Orm) transformUpdateSql() string {
	return "UPDATE " + v.transformTable() + " " + v.transformQuery()
}

func (v *Orm) transformDeleteSql() string {
	return "DELETE " + v.transformTable() + " " + v.transformQuery()
}

func (v *Orm) transformInsertSql() string {
	return "INSERT INTO " + v.transformTable() + " " + v.transformQuery()
}

func (v *Orm) Err() error {
	return v.err
}

func (v *Orm) setErr(e error) *Orm {
	v.err = e
	if v.db.errHandler != nil {
		v.db.errHandler(e)
	}
	return v
}

func (v *Orm) transformTable() string {
	if len(v.table) > 1 {
		v.table.With(mx.WithTable)
	}
	v.table.With(mx.WithBackquote)
	return v.table.GetQuery()
}
func (v *Orm) transformQuery() string {
	if len(v.table) > 1 {
		v.mix.With(mx.WithTable)
	}
	v.mix.With(mx.WithBackquote)
	return v.mix.GetQuery()
}

func (v *Orm) GetArgs() []interface{} {
	if len(v.table) > 1 {
		v.mix.With(mx.WithTable)
	}
	v.mix.With(mx.WithBackquote)
	return v.mix.GetArgs()
}
