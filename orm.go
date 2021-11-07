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

	query          string
	args           []interface{}
	dest           interface{}
	err            error
	rawFields      bool
	wk             []string
	wv             []interface{}
	sk             []string
	sv             []interface{}
	StartQueryTime time.Time
	Sql            string
}

func (v *Orm) Query(query string, args ...interface{}) *Orm {
	v.addMix(mix.NewMix(query, args...))
	return v
}

func (v *Orm) addMix(m mx.Mix, typs ...string) *Orm {

	if v.mix == nil {
		v.mix = make(mx.Mixs, 0)
	}

	if len(typs) > 0 && v.mixType != typs[0] {
		typ := typs[0]
		if typ == "where" {
			v.mix = append(v.mix, mix.NewMix(" WHERE "))
		} else if typ == "set" {
			v.mix = append(v.mix, mix.NewMix(" SET "))
		} else if typ == "group" {
			v.mix = append(v.mix, mix.NewMix(" GROUP BY "))
		} else if typ == "limit" {
			v.mix = append(v.mix, mix.NewMix(" LIMIT "))
		} else if typ == "order" {
			v.mix = append(v.mix, mix.NewMix(" ORDER BY "))
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

func (v *Orm) Field(fields ...interface{}) *Orm {
	for _, f := range fields {
		k := field.GetField(f)
		if k != nil {
			v.addField(k)
		}
	}
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
	return "SELECT " + v.transformFields() + " FROM " + v.table.GetQuery() + " " + v.transformQuery()
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

func (v *Orm) WhereStru(s interface{}) *Orm {
	rv := stringify.GetReflectValue(s)

	p := mx.ConditionMix{}

	loopStruct(rv, func(v reflect.Value, s reflect.StructField) bool {
		db := s.Tag.Get("db")
		dbset := s.Tag.Get("dbwhere")
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
			p = append(p, mix.NewMix("%t=?", field.NewField(db), v.Interface()))

			return true
		}
		return false
	})

	return v.addMix(p, "where")
}

func (v *Orm) WhereMap(s ...interface{}) *Orm {

	if len(s)%2 == 1 {
		v.setErr(ErrOddNumberOfParams)
		return v
	}

	if v.wk == nil {
		v.wk = make([]string, 0)
	}

	if v.wv == nil {
		v.wv = make([]interface{}, 0)
	}

	for k, vs := range s {
		if k%2 == 0 {
			v.wk = append(v.wk, stringify.ToString(vs))
		} else {
			v.wv = append(v.wv, vs)
		}
	}
	return v
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
	p := map[string]interface{}{}
	rv := stringify.GetReflectValue(s)
	loopStruct(rv, func(v reflect.Value, s reflect.StructField) bool {
		db := s.Tag.Get("db")
		dbset := s.Tag.Get("dbset")
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
			p[db] = v.Interface()
			return true
		}
		return false
	})
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

func (v *Orm) SetMap(s ...interface{}) *Orm {
	p, err := sliceToMap(s)
	if err != nil {
		v.setErr(err)
	}
	return v.Set(p)
}

func (v *Orm) transformQuery() string {
	return v.mix.GetQuery()
}

func (v *Orm) GetArgs() []interface{} {
	return v.args
}
