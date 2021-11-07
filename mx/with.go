package mx

type With int

const (
	WithNone With = iota
	WithAs
	WithTable
	WithBackquote
)

type with interface {
	With(With)
}

type WithTrait struct {
	withAs        bool
	withBackquote bool
	withTable     bool
	query         bool
}

func (wt *WithTrait) With(w With) {
	if w == WithAs {
		wt.withAs = true
	}
	if w == WithBackquote {
		wt.withBackquote = true
	}
	if w == WithTable {
		wt.withTable = true
	}
}

func (wt *WithTrait) SetQuery() {
	wt.query = true
}

func (wt *WithTrait) Reset() {
	wt.withAs = false
	wt.withBackquote = false
	wt.query = false
	wt.withTable = false
}

func (wt *WithTrait) IsWithAs() bool {
	return wt.withAs
}

func (wt *WithTrait) IsWithBackquote() bool {
	return wt.withBackquote
}

func (wt *WithTrait) IsWithTable() bool {
	return wt.withTable
}
func (wt *WithTrait) IsQuery() bool {
	return wt.query
}
