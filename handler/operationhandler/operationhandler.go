package operationhandler

import (
	"rfs/models/entity"
	"rfs/models/modelconst"
)

type OperationHandler struct {
	operations    []*entity.Operation
	OperationChan chan *entity.Operation
}

func NewOperationHandler() *OperationHandler {
	operationHandler := &OperationHandler{
		operations:    make([]*entity.Operation, 0),
		OperationChan: make(chan *entity.Operation, 1),
	}

	go operationHandler.listenOperationChannel()

	return operationHandler
}

func (OperationHandler *OperationHandler) GetNewOperations() []*entity.Operation {
	var newOperations []*entity.Operation

	for _, op := range OperationHandler.operations {
		if op.State == modelconst.NEW {
			op.State = modelconst.PENDING
			newOperations = append(newOperations, op)
		}
	}

	return newOperations
}

func (operationhandler *OperationHandler) listenOperationChannel() {
	for op := range operationhandler.OperationChan {
		operationhandler.operations = append(operationhandler.operations, op)
	}
}
