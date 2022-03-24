package message

import "rfs/models/entity"

type OperationChanCommand int8

const (
	ADD OperationChanCommand = iota + 1
	REMOVE
)

type OperationChanMsg struct {
	Operation *entity.Operation
	Command   OperationChanCommand
}

func NewOperationMsg(operation *entity.Operation, command OperationChanCommand) *OperationChanMsg {
	return &OperationChanMsg{
		Operation: operation,
		Command:   command,
	}
}
