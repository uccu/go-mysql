package table

import "github.com/uccu/go-mysql/mx"

type join struct {
	table         mx.Container
	joinType      mx.JoinType
	joinCondition mx.ConditionMix
}

type joins []*join
