package modelconst

import "strconv"

type OperationType int8

const (
	CREATE_FILE OperationType = iota + 1
	APPEND_RECORD
)

func (opType OperationType) String() string {
	return strconv.Itoa(int(opType))
}
