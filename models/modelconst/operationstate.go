package modelconst

type OperationState int8

const (
	NEW OperationState = iota + 1
	PENDING
	CONFIRMED
)
