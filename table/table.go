package table

import (
	"github.com/uccu/go-mysql/mx"
)

type join struct {
	table         mx.Container
	joinType      mx.JoinType
	joinCondition mx.ConditionMix
}

type joins []*join

type Table struct {
	Name    string
	Alias   string
	RawName string
	suffix  func(interface{}) string
	join    joins

	with mx.WithTrait
}

func (t *Table) With(w mx.With) {
	t.with.With(w)
}

func (t *Table) GetName() string {
	alias := t.Alias
	if t.Alias == "" {
		alias = t.RawName
	}
	if t.with.IsWithBackquote() {
		alias = "`" + alias + "`"
	}
	if !t.with.IsQuery() {
		t.with.Reset()
	}
	return alias
}
func (t *Table) Suffix(s interface{}) {
	t.RawName += t.suffix(s)
}

func (t *Table) GetQuery() string {

	t.with.SetQuery()
	query := t.RawName

	if t.with.IsWithBackquote() {
		query = "`" + query + "`"
	}

	if t.with.IsWithAlias() && t.Alias != "" {
		query += " " + t.GetName()
	}

	if t.join != nil {
		for _, j := range t.join {
			if t.with.IsWithAlias() {
				j.table.With(mx.WithAlias)
			}
			if t.with.IsWithBackquote() {
				j.table.With(mx.WithBackquote)
			}

			query += " " + j.joinType.String() + " " + j.table.GetQuery() + " " + j.joinCondition.GetQuery()
		}
	}
	t.with.Reset()

	return query
}

func (t *Table) GetArgs() []interface{} {

	t.with.SetQuery()
	args := []interface{}{}
	if t.join != nil {
		for _, j := range t.join {
			args = append(args, j.joinCondition.GetArgs()...)
		}
	}
	if len(args) == 0 {
		return nil
	}
	t.with.Reset()
	return args
}

func NewTable(name, prefix string, suffix func(interface{}) string) *Table {
	return &Table{Name: name, RawName: prefix + name, suffix: suffix}
}
