package field

import (
	"github.com/uccu/go-mysql/mix"
	"github.com/uccu/go-mysql/mx"
)

type RawField struct {
	q string
}

func (f *RawField) GetQuery() string {
	return f.q
}

func (t *RawField) With(w mx.With) mx.Field {
	return t
}

func (f *RawField) GetArgs() []interface{} {
	return nil
}

func NewRawField(q string) *RawField {
	return &RawField{q}
}

func (f *RawField) ToMix() mx.Mix {
	return &mix.Field{Field: f}
}
