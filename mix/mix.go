package mix

import (
	"strings"

	"github.com/uccu/go-mysql/mx"
)

type Mix struct {
	q string
	m mx.Mixs
	a []interface{}
}

func (m *Mix) GetQuery() string {
	query := m.q
	for _, m := range m.m {
		query = strings.Replace(query, "%t", m.GetQuery(), 1)
	}
	return query
}

func (m *Mix) GetArgs() []interface{} {
	return m.a
}

func NewMix(q string, f ...interface{}) *Mix {

	mix := &Mix{q: q}

	mixs := mx.Mixs{}
	args := []interface{}{}
	for _, f := range f {
		if f, ok := f.(mx.Mix); ok {
			mixs = append(mixs, f)
			args = append(args, f.GetArgs()...)
			continue
		}
		args = append(args, f)
	}

	if len(mixs) > 0 {
		mix.m = mixs
	}

	if len(args) > 0 {
		mix.a = args
	}

	return mix
}
