package minehandler

import (
	"rfs/handler/chainhandler"
	"rfs/handler/operationhandler"
	"rfs/models/entity"
	"rfs/secsuit"
	"time"
)

type MinerHandler struct {
	genOpBlockTimeout  time.Duration
	genNoOpBlockTime   time.Time
	cancelNoOpBlockGen chan int
	// blockchain         *entity.BlockChain

	operationHandler *operationhandler.OperationHandler
	chainhandler     *chainhandler.ChainHandler
}

func NewMinerHandler() *MinerHandler {

	return &MinerHandler{
		genOpBlockTimeout:  3 * time.Second,
		genNoOpBlockTime:   time.Now(),
		cancelNoOpBlockGen: make(chan int),

		operationHandler: operationhandler.NewOperationHandler(),
		chainhandler:     chainhandler.NewChainHandler(),
	}
}

func (minerHandler *MinerHandler) AddNewOperation(operation *entity.Operation) {
	minerHandler.operationHandler.OperationChan <- operation
}

func (minerHandler *MinerHandler) GenerateOpBlock() {

	time.Sleep(3 * time.Second)

	newOperations := minerHandler.operationHandler.GetNewOperations()

	lastblock := minerHandler.chainhandler.GetLongestValidChain()

	newOpBlock := &entity.Block{
		PrevHash:   secsuit.ComputeHash(lastblock.String()),
		Operations: newOperations,
		TimeStamp:  time.Now(),
	}

	minerHandler.chainhandler.AddBlock(newOpBlock)

	// for {
	// 	time.After(worker.genOpBlockTimeout)
	// 	if len(worker.operations) == 0 {
	// 		continue
	// 	}
	// 	worker.cancelNoOpBlockGen <- 1

	// 	time.Sleep(time.Second)

	// }
}

func (minerHandler *MinerHandler) GenerateNoOpBlock() {

	time.Sleep(time.Second)

	lastblock := minerHandler.chainhandler.GetLongestValidChain()

	newOpBlock := &entity.Block{PrevHash: secsuit.ComputeHash(lastblock.String())}

	minerHandler.chainhandler.AddBlock(newOpBlock)

	// for {

	// }
}
