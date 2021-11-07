package field

import "github.com/uccu/go-mysql/mx"

type RawField struct {
	q string
}

func (f *RawField) GetQuery() string {
	return f.q
}

func (t *RawField) With(w mx.With) mx.Field {
	return t
}

func NewRawField(q string) *RawField {
	return &RawField{q}
}

func (f *RawField) GetArgs() []interface{} {
	return nil
}
