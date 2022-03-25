package operationhandler

import (
	"cfs/models/entity"
	"cfs/models/message"
	"cfs/models/modelconst"
	"cfs/sharedchannel"
	"fmt"
	"log"
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
	newOperations := OperationHandler.operations

	OperationHandler.operations = make([]*entity.Operation, 0)

	return newOperations
}

func (OperationHandler *OperationHandler) SetOperationsStatus(operations []*entity.Operation, opState modelconst.OperationState) {

	// for _, op := range operations {
	// 	if op.State == modelconst.NEW {
	// 		op.State = opState
	// 	}
	// }
}

//Todo: check and resolve concurrancy
func (operationHandler *OperationHandler) RemoveOperations(operations []*entity.Operation) {
	for _, op := range operations {
		operationHandler.sharedchannel.InternalOperationChan <- message.NewOperationMsg(op, message.REMOVE)
	}
}

func (operationhandler *OperationHandler) ListenOperationChannel() {
	for op := range operationhandler.sharedchannel.InternalOperationChan {
		if op.Command == message.ADD {
			if !operationhandler.validateOperation(op.Operation) {
				log.Println("OperationHandler/ListenOperationChannel - invalid operation ", op)
				continue
			}

			operationhandler.sharedchannel.BroadcastOperationChan <- op.Operation

			operationhandler.operations = append(operationhandler.operations, op.Operation)

		} else {
			for index, operation := range operationhandler.operations {
				if op.Operation.OperationId == operation.OperationId {
					operationhandler.operations = append(operationhandler.operations[:index], operationhandler.operations[index+1:]...)
				}
			}
		}
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
