package mx

type Query interface {
	GetQuery() string
}

type Args interface {
	GetArgs() []interface{}
}

type Value interface {
	Query
}
