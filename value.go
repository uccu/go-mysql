package mysql

import "strings"

type mutiValue struct {
	q string
	f []value
}

func (f *mutiValue) GetQuery() string {
	query := f.q
	for _, f := range f.f {
		query = strings.Replace(query, "?", f.GetQuery(), 1)
	}
	return query
}

func MutiValue(query string, v ...value) *mutiValue {
	return &mutiValue{query, v}
}
