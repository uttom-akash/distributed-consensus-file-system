package operationhandler

import (
	"fmt"
	"log"
	"rfs/models/entity"
	"rfs/models/modelconst"
	"rfs/sharedchannel"
	"sync"
)

type OperationHandler struct {
	operations    []*entity.Operation
	sharedchannel *sharedchannel.SharedChannel
}

func NewOperationHandler() IOperationHandler {
	operationHandler := &OperationHandler{
		operations:    make([]*entity.Operation, 0),
		sharedchannel: sharedchannel.NewSingletonSharedChannel(),
	}

	return operationHandler
}

var lock = &sync.Mutex{}
var singletonInstance IOperationHandler

func NewSingletonOperationHandler() IOperationHandler {

	if singletonInstance == nil {
		lock.Lock()
		defer lock.Unlock()

		if singletonInstance == nil {
			fmt.Println("Creating single instance now.")
			singletonInstance = NewOperationHandler()
		} else {
			fmt.Println("Single instance already created.")
		}
	} else {
		fmt.Println("Single instance already created.")
	}

	return singletonInstance
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

func (OperationHandler *OperationHandler) SetOperationsStatus(operations []*entity.Operation, opState modelconst.OperationState) {

	for _, op := range operations {
		if op.State == modelconst.NEW {
			op.State = opState
		}
	}
}

func (operationhandler *OperationHandler) ListenOperationChannel() {
	for op := range operationhandler.sharedchannel.Operation {
		if !operationhandler.validateOperation(op) {
			log.Println("OperationHandler/ListenOperationChannel - invalid operation ", op)
			continue
		}

		operationhandler.sharedchannel.BroadcastOperation <- op

		operationhandler.operations = append(operationhandler.operations, op)

	}
}

func (operationhandler *OperationHandler) validateOperation(operation *entity.Operation) bool {

	for _, op := range operationhandler.operations {
		if op.OperationId == operation.OperationId {
			return false
		}
	}

	return true
}
