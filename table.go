package mysql

import "github.com/uccu/go-stringify"

type JoinType byte

const (
	NO_JOIN JoinType = iota
	JOIN
	LEFT_JOIN
	RIGHT_JOIN
	INNER_JOIN
)

type table struct {
	Name          string
	As            string
	RawName       string
	db            *DB
	join          JoinType
	joinCondition Condition
}

type Table interface {
	GetAs(...With) string
	GetQuery() string
}

type Tables []Table

func (ts Tables) GetQuery() string {
	s := []string{}
	for _, t := range ts {
		s = append(s, t.GetQuery())
	}
	return stringify.ToString(s, ", ")
}

func (t *table) GetAs(w ...With) string {
	as := t.As
	if t.As == "" {
		as = t.RawName
	}
	if len(w) > 0 && w[0]&WithBackquote > 0 {
		as = "`" + as + "`"
	}
	return t.As
}
func (t *table) Suffix(s interface{}) {
	t.RawName += t.db.suffix(s)
}

func (t *table) GetQuery() string {
	query := "`" + t.RawName + "`"
	if t.As != "" {
		query += " AS " + "`" + t.As + "`"
	}
	if t.join != NO_JOIN {
		query += " ON " + t.joinCondition.GetQuery()
	}
	return query
}

func (o *Orm) Suffix(s interface{}) *Orm {
	for _, t := range o.tables {
		if t, ok := t.(*table); ok {
			t.Suffix(s)
		}
	}
	return o
}
