package operationhandler

import "rfs/models/entity"

type IOperationHandler interface {
	GetNewOperations() []*entity.Operation
	SetOperationsPending(operations []*entity.Operation)
	ListenOperationChannel()
}
