package field

import (
	"github.com/uccu/go-mysql/mix"
	"github.com/uccu/go-mysql/mx"
)

type Field struct {
	Table string
	Name  string
	Alias string

	with mx.WithTrait
}

func NewField(name string) *Field {
	return &Field{Name: name}
}

func (t *Field) With(w mx.With) mx.Field {
	t.with.With(w)
	return t
}

func (f *Field) SetAlias(n string) *Field {
	f.Alias = n
	return f
}

func (f *Field) SetTable(n string) *Field {
	f.Table = n
	return f
}

func (t *Field) GetName() string {
	alias := t.Alias
	if t.Alias == "" {
		alias = t.Name
	}
	if t.with.IsWithBackquote() {
		alias = "`" + alias + "`"
	}
	if !t.with.IsQuery() {
		t.with.Reset()
	}
	return alias
}

func (f *Field) GetQuery() string {

	f.with.SetQuery()
	query := f.Name

	if f.with.IsWithBackquote() {
		query = "`" + query + "`"
	}

	if f.with.IsWithAlias() && f.Alias != "" {
		query += " " + f.GetName()
	}

	if f.with.IsWithTable() && f.Table != "" {
		tableName := f.Table
		if f.with.IsWithBackquote() {
			tableName = "`" + tableName + "`"
		}
		query = tableName + "." + query
	}
	f.with.Reset()
	return query
}

func (f *Field) GetArgs() []interface{} {
	return nil
}

func (f *Field) ToMix() mx.Mix {
	return &mix.Field{Field: f}
}
