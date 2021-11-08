package mix

import "github.com/uccu/go-mysql/mx"

type Field struct {
	mx.Field
}

func (t *Field) With(w mx.With) mx.Mix {
	t.Field.With(w)
	return t
}
