package mysql

import "strings"

type value interface {
	GetQuery() string
}

type MutiValue struct {
	q string
	f []value
}

func (f *MutiValue) GetQuery() string {
	query := f.q
	for _, f := range f.f {
		query = strings.Replace(query, "%t", f.GetQuery(), 1)
	}
	return query
}

func NewMutiValue(query string, v ...value) *MutiValue {
	return &MutiValue{query, v}
}
