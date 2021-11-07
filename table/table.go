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
	As      string
	RawName string
	suffix  func(interface{}) string
	join    joins

	with mx.WithTrait
}

func (t *Table) With(w mx.With) {
	t.with.With(w)
}

func (t *Table) GetName() string {
	as := t.As
	if t.As == "" {
		as = t.RawName
	}
	if t.with.IsWithBackquote() {
		as = "`" + as + "`"
	}
	if !t.with.IsQuery() {
		t.with.Reset()
	}
	return as
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

	if t.with.IsWithAs() && t.As != "" {
		query += " " + t.GetName()
	}

	if t.join != nil {
		for _, j := range t.join {
			if t.with.IsWithAs() {
				j.table.With(mx.WithAs)
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
