package mx

import "github.com/uccu/go-stringify"

type Field interface {
	Query
	Args
	With(With) Field
	ToMix() Mix
}

type Fields []Field

func (f Fields) GetQuery() string {

	if len(f) == 0 {
		return ""
	}
	list := []string{}
	for _, f := range f {
		list = append(list, f.GetQuery())
	}

	return stringify.ToString(list, ", ")
}

func (f Fields) With(w With) {
	for _, f := range f {
		f.With(w)
	}
}
