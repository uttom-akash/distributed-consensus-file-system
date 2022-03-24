package operationhandler

import (
	"rfs/models/entity"
	"rfs/models/modelconst"
)

type IOperationHandler interface {
	GetNewOperations() []*entity.Operation
	SetOperationsStatus(operations []*entity.Operation, opState modelconst.OperationState)
	RemoveOperations(operations []*entity.Operation)
	ListenOperationChannel()
}
