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

type Container interface {
	Query
	Args
	with
}
