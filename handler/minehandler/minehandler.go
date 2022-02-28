package minehandler

import (
	"context"
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

	for {

		time.Sleep(minerHandler.genOpBlockTimeout)

		newOperations := minerHandler.operationHandler.GetNewOperations()

		if len(newOperations) == 0 {
			continue
		}

		minerHandler.cancelNoOpBlockGen <- 1

		lastblock := minerHandler.chainhandler.GetLongestValidChain()

		newOpBlock := &entity.Block{
			PrevHash:   secsuit.ComputeHash(lastblock.String()),
			Operations: newOperations,
			TimeStamp:  time.Now(),
		}

		minerHandler.chainhandler.AddBlock(newOpBlock)

		minerHandler.cancelNoOpBlockGen <- 2
	}
}

func (minerHandler *MinerHandler) GenerateNoOpBlock1() {
	contxt, cancelFunc := context.WithCancel(context.TODO())
	minerHandler.GenerateNoOpBlock(contxt)

	for {
		msg := <-minerHandler.cancelNoOpBlockGen
		if msg == 1 {
			cancelFunc()
		}

		if msg == 2 {
			contxt, cancelFunc = context.WithCancel(context.TODO())
			minerHandler.GenerateNoOpBlock()
		}
	}
}

func (minerHandler *MinerHandler) GenerateNoOpBlock(ctx context.Context) {

	for {
		time.Sleep(time.Second)

		lastblock := minerHandler.chainhandler.GetLongestValidChain()

		newNoOpBlock := &entity.Block{PrevHash: secsuit.ComputeHash(lastblock.String())}

		select {
		case <-ctx.Done():
			break
		case:
			minerHandler.chainhandler.AddBlock(newNoOpBlock)
		}
	}
}
