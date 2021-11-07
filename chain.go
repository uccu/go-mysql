package mysql

import (
	"reflect"

	"github.com/uccu/go-mysql/field"
	"github.com/uccu/go-mysql/mix"
	"github.com/uccu/go-mysql/mx"
	"github.com/uccu/go-stringify"
)

func (v *Orm) Field(fields ...interface{}) *Orm {
	for _, f := range fields {
		k := field.GetField(f)
		if k != nil {
			v.addField(k)
		}
	}
	return v
}

func (v *Orm) whereStru(s interface{}) *Orm {
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

func (v *Orm) Where(s ...interface{}) *Orm {

	if len(s) == 1 {
		rv := stringify.GetReflectValue(s[0])
		if rv.Kind() == reflect.Struct {
			return v.whereStru(s[0])
		} else if rv.Kind() == reflect.Map {
			return v.whereMap(s[0].(map[string]interface{}))
		}
	}

	if len(s)%2 == 1 {
		v.setErr(ErrOddNumberOfParams)
		return v
	}

	p := mx.ConditionMix{}
	for k := 0; k < len(s); k += 2 {
		p = append(p, mix.NewMix("%t=?", field.NewField(s[k].(string)), s[k+1]))
	}

	return v.addMix(p, "where")
}

func (v *Orm) whereMap(data map[string]interface{}) *Orm {

	p := mx.ConditionMix{}
	for k, v := range data {
		p = append(p, mix.NewMix("%t=?", field.NewField(k), v))
	}

	return v.addMix(p, "where")
}

func (v *Orm) setStru(s interface{}) *Orm {

	rv := stringify.GetReflectValue(s)

	p := mx.SliceMix{}

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
			p = append(p, mix.NewMix("%t=?", field.NewField(db), v.Interface()))
			return true
		}
		return false
	})
	return v.addMix(p, "set")
}

func (v *Orm) Set(s ...interface{}) *Orm {

	if len(s) == 1 {
		rv := stringify.GetReflectValue(s[0])
		if rv.Kind() == reflect.Struct {
			return v.setStru(s[0])
		} else if rv.Kind() == reflect.Map {
			return v.setMap(s[0].(map[string]interface{}))
		}
	}

	if len(s)%2 == 1 {
		v.setErr(ErrOddNumberOfParams)
		return v
	}

	p := mx.SliceMix{}
	for k := 0; k < len(s); k += 2 {
		p = append(p, mix.NewMix("%t=?", field.NewField(s[k].(string)), s[k+1]))
	}

	return v.addMix(p, "set")
}

func (v *Orm) setMap(data map[string]interface{}) *Orm {

	p := mx.SliceMix{}
	for k, v := range data {
		p = append(p, mix.NewMix("%t=?", field.NewField(k), v))
	}

	return v.addMix(p, "set")
}
