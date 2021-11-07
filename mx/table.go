package mx

import "github.com/uccu/go-stringify"

type Table interface {
	GetName() string
	Container
}

type Tables []Table

func (ts Tables) With(w With) Tables {
	for _, t := range ts {
		t.With(w)
	}
	return ts
}

func (ts Tables) GetQuery() string {
	s := []string{}
	for _, t := range ts {
		s = append(s, t.GetQuery())
	}
	return stringify.ToString(s, ", ")
}

func (ts Tables) GetArgs() []interface{} {
	args := []interface{}{}
	for _, t := range ts {
		args = append(args, t.GetArgs()...)
	}
	if len(args) == 0 {
		return nil
	}
	return args
}
