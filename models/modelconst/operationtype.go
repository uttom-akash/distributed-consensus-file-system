package modelconst

type OperationType int8

const (
	CREATE_FILE OperationType = iota + 1
	APPEND_RECORD
)
