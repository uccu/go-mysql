package mix

import (
	"github.com/uccu/go-mysql/mx"
)

type Raw struct {
	q string
}

func (t *Raw) With(w mx.With) {
}

func (m *Raw) GetQuery() string {
	return m.q
}

func (m *Raw) GetArgs() []interface{} {
	return nil
}

func NewRawMix(q string) *Raw {
	return &Raw{q}
}
