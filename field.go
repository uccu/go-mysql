package mysql

import "strings"

type Field struct {
	t    Table
	Name string
	As   string
}

type With int

const (
	WithAs With = 1 << iota
	WithTable
	WithBackquote
)

func (t *Field) GetAs(w ...With) string {
	as := t.As
	if t.As == "" {
		as = t.Name
	}
	if len(w) > 0 && w[0]&WithBackquote > 0 {
		as = "`" + as + "`"
	}
	return t.As
}

func (f *Field) GetQuery(w ...With) string {

	query := f.Name
	if len(w) == 0 {
		return query
	}

	if w[0]&WithAs > 0 && f.As != "" {
		query += " AS " + f.GetAs(w...)
	}

	if w[0]&WithTable > 0 {
		query = f.t.GetAs(w...) + "." + query
	}

	return ""
}

type mutiField struct {
	q string
	f []*Field
}

func (f *mutiField) GetQuery() string {
	query := f.q
	for _, f := range f.f {
		query = strings.Replace(query, "%t", f.GetQuery(), 1)
	}
	return query
}

func MutiField(query string, fields ...*Field) *mutiField {
	return &mutiField{query, fields}
}

type rawField struct {
	q string
}

func (f *rawField) GetQuery() string {
	return f.q
}

func RawField(query string) *rawField {
	return &rawField{query}
}
