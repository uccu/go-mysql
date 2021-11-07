package field

import (
	"strings"

	"github.com/uccu/go-mysql/mx"
)

type MutiField struct {
	q string
	f mx.Fields
}

func (f *MutiField) GetQuery() string {
	query := f.q
	for _, f := range f.f {
		query = strings.Replace(query, "?", f.GetQuery(), 1)
	}
	return query
}

func NewMutiField(q string, f ...mx.Field) *MutiField {
	return &MutiField{q, f}
}

func (f *MutiField) GetArgs() []interface{} {
	return nil
}
