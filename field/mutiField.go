package field

import (
	"strings"

	"github.com/uccu/go-mysql/mix"
	"github.com/uccu/go-mysql/mx"
)

type MutiField struct {
	q string
	f mx.Fields
}

func (t *MutiField) With(w mx.With) mx.Field {
	for _, f := range t.f {
		f.With(w)
	}
	return t
}

func (m *MutiField) GetQuery() string {
	query := m.q
	for _, f := range m.f {
		query = strings.Replace(query, "%t", f.GetQuery(), 1)
	}
	return query
}

func NewMutiField(q string, f ...mx.Field) *MutiField {
	return &MutiField{q: q, f: f}
}

func (f *MutiField) GetArgs() []interface{} {
	return nil
}

func (f *MutiField) ToMix() mx.Mix {
	return &mix.Field{Field: f}
}
