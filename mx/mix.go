package mx

import "github.com/uccu/go-stringify"

type Mix interface {
	Query
	Args
	With(With) Mix
}

type Mixs []Mix

func (f Mixs) With(w With) Mix {
	for _, f := range f {
		f.With(w)
	}
	return f
}

func (f Mixs) GetQuery() string {
	if len(f) == 0 {
		return ""
	}
	list := []string{}
	for _, f := range f {
		list = append(list, f.GetQuery())
	}

	return stringify.ToString(list, "")
}

func (f Mixs) GetArgs() []interface{} {
	a := []interface{}{}
	for _, f := range f {
		a = append(a, f.GetArgs()...)
	}
	if len(a) == 0 {
		return nil
	}
	return a
}

type ConditionMix []Mix

func (f ConditionMix) GetQuery() string {
	if len(f) == 0 {
		return ""
	}
	list := []string{}
	for _, f := range f {
		list = append(list, f.GetQuery())
	}

	return stringify.ToString(list, " AND ")
}

func (f ConditionMix) GetArgs() []interface{} {
	return Mixs(f).GetArgs()
}

func (f ConditionMix) With(w With) Mix {
	return Mixs(f).With(w)
}

type SliceMix []Mix

func (f SliceMix) GetQuery() string {
	if len(f) == 0 {
		return ""
	}
	list := []string{}
	for _, f := range f {
		list = append(list, f.GetQuery())
	}

	return stringify.ToString(list, ", ")
}

func (f SliceMix) GetArgs() []interface{} {
	return Mixs(f).GetArgs()
}

func (f SliceMix) With(w With) Mix {
	return Mixs(f).With(w)
}
