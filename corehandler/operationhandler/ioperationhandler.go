package operationhandler

import (
	"cfs/models/entity"
	"cfs/models/modelconst"
)

type IOperationHandler interface {
	GetNewOperations() []*entity.Operation
	SetOperationsStatus(operations []*entity.Operation, opState modelconst.OperationState)
	RemoveOperations(operations []*entity.Operation)
	ListenOperationChannel()
}
