package operationhandler

import (
	"fmt"
	"rfs/models/entity"
	"rfs/models/modelconst"
	"sync"
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

	return operationHandler
}

var lock = &sync.Mutex{}
var singletonInstance *OperationHandler

func NewSingletonOperationHandler() *OperationHandler {

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

func (operationhandler *OperationHandler) ListenOperationChannel() {
	for op := range operationhandler.OperationChan {
		operationhandler.operations = append(operationhandler.operations, op)
	}
}
