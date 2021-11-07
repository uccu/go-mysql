package mix

import (
	"strings"

	"github.com/uccu/go-mysql/mx"
)

type Mix struct {
	q string
	m mx.Mixs
	a []interface{}

	with mx.WithTrait
}

func (t *Mix) With(w mx.With) mx.Mix {
	t.with.With(w)
	return t
}

func (m *Mix) GetQuery() string {
	m.with.SetQuery()
	query := m.q
	for _, v := range m.m {
		query = strings.Replace(query, "%t", v.With(m.with.Status).GetQuery(), 1)
	}
	m.with.Reset()
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
