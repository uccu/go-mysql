package field

import (
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

func (f *Field) SetAs(n string) {
	f.Alias = n
}

func (f *Field) SetTable(n string) {
	f.Table = n
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

func GetField(f interface{}) mx.Field {
	if k, ok := f.(mx.Field); ok {
		return k
	} else if k, ok := f.(string); ok {
		return NewField(k)
	}
	return nil
}

func (f *Field) GetArgs() []interface{} {
	return nil
}
