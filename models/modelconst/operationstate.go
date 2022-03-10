package modelconst

import "strconv"

type OperationState int8

const (
	NEW OperationState = iota + 1
	PENDING
	CONFIRMED
)

func (opState OperationState) String() string {
	return strconv.Itoa(int(opState))
}
