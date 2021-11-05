package mysql

import (
	"strings"

	"github.com/uccu/go-stringify"
)

type Condition interface {
	GetQuery() string
	GetArgs() []interface{}
}

type Conditions []Condition

func (f Conditions) GetQuery() string {

	if len(f) == 0 {
		return ""
	}

	list := []string{}
	for i := 0; i < len(f); i++ {
		list = append(list, "(%t)")
	}

	return MutiCondition(stringify.ToString(list, " AND "), f...).GetQuery()
}

func (f Conditions) GetArgs() []interface{} {
	a := []interface{}{}
	for _, f := range f {
		a = append(a, f.GetArgs()...)
	}
	if len(a) == 0 {
		return nil
	}
	return a
}

type mutiCondition struct {
	q string
	f Conditions
}

func (f *mutiCondition) GetQuery() string {
	query := f.q
	for _, f := range f.f {
		query = strings.Replace(query, "%t", f.GetQuery(), 1)
	}
	return query
}

func (f *mutiCondition) GetArgs() []interface{} {
	return f.GetArgs()
}

func MutiCondition(query string, c ...Condition) *mutiCondition {
	return &mutiCondition{query, c}
}

type rawCondition struct {
	q string
	a []interface{}
}

func RawCondition(q string, a ...interface{}) *rawCondition {
	return &rawCondition{q, a}
}

func (f *rawCondition) GetQuery() string {
	return f.q
}

func (f *rawCondition) GetArgs() []interface{} {
	return f.a
}
