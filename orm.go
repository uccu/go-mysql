package mysql

import (
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
	subTable  bool
	subValue  int64
}

func (v *Orm) SubValue(s int64) *Orm {
	v.subTable = true
	v.subValue = s
	return v
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

func (v *Orm) transformFields() string {
	if len(v.fields) == 0 {
		if v.dest != nil {
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
					v.setErr(NO_DB_TAG)
					return "1"
				}

				return v.Fields(removeRep(fields)).transformFields()
			}
		}
		return "*"
	}
	if v.rawFields {
		return strings.Join(v.fields, ",")
	}
	return "`" + strings.Join(v.fields, "`,`") + "`"
}

func (v *Orm) transformTableName() string {
	table := v.table
	if v.subTable {
		table = v.db.subTable(v.table, v.subValue)
	}
	return "`" + v.db.prefix + table + "`"
}

func (v *Orm) transformSelectSql() string {
	return "SELECT " + v.transformFields() + " FROM " + v.transformTableName() + v.transformQuery()
}

func (v *Orm) RawFields(r bool) *Orm {
	v.rawFields = r
	return v
}

func (v *Orm) Err() error {
	return v.err
}

func (v *Orm) setErr(e error) *Orm {
	v.err = e
	if v.db.errHandle != nil {
		v.db.errHandle(e)
	}
	return v
}

func (v *Orm) WhereStru(s interface{}) *Orm {
	p := map[string]interface{}{}
	rv := stringify.GetReflectValue(s)
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
		if v.CanAddr() && v.CanInterface() {
			p[db] = v.Addr().Interface()
			return true
		}
		return false
	})
	return v.Where(p)
}

func (v *Orm) WhereMap(s ...interface{}) *Orm {

	p := map[string]interface{}{}
	if len(s)%2 == 1 {
		v.setErr(ODD_PARAM)
	} else {
		for k := range s {
			if k%2 == 1 {
				p[stringify.ToString(s[k-1])] = s[k]
			}
		}
	}
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
		if v.CanAddr() && v.CanInterface() {
			p[db] = v.Addr().Interface()
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

func (v *Orm) transformQuery() string {

	where := ""
	if len(v.wk) > 0 {
		v.args = append(v.wv, v.args...)
		where = " WHERE `" + strings.Join(v.wk, "`=? AND `") + "`=?"
	}

	set := ""
	if len(v.sk) > 0 {
		v.args = append(v.sv, v.args...)
		set = " SET `" + strings.Join(v.sk, "`=?,`") + "`=?"
	}

	query := ""
	if v.query != "" {
		query = " " + v.query
	}

	return set + where + query
}
