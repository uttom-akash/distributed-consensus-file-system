package minehandler

import (
	"fmt"
	"rfs/handler/chainhandler"
	"rfs/handler/operationhandler"
	"rfs/models/entity"
	"rfs/sharedchannel"
	"sync"
	"time"
)

type MinerHandler struct {
	genOpBlockTimeout  time.Duration
	genNoOpBlockTime   time.Time
	cancelNoOpBlockGen chan int

	operationHandler *operationhandler.OperationHandler
	chainhandler     *chainhandler.ChainHandler
	sharedchannel    *sharedchannel.SharedChannel
}

func NewMinerHandler() *MinerHandler {

	minerHandler := &MinerHandler{
		genOpBlockTimeout:  2 * time.Minute,
		genNoOpBlockTime:   time.Now(),
		cancelNoOpBlockGen: make(chan int),
		sharedchannel:      sharedchannel.NewSingletonSharedChannel(),
		operationHandler:   operationhandler.NewSingletonOperationHandler(),
		chainhandler:       chainhandler.NewSingletonChainHandler(),
	}

	return minerHandler
}

var lock = &sync.Mutex{}
var singletonInstance *MinerHandler

func NewSingletonMinerHandler() *MinerHandler {

	if singletonInstance == nil {
		lock.Lock()
		defer lock.Unlock()

		if singletonInstance == nil {
			fmt.Println("Creating single instance now.")
			singletonInstance = NewMinerHandler()
		} else {
			fmt.Println("Single instance already created.")
		}
	} else {
		fmt.Println("Single instance already created.")
	}

	return singletonInstance
}

func (minerHandler *MinerHandler) AddNewOperation(operation *entity.Operation) {
	minerHandler.sharedchannel.Operation <- operation
}

func (minerHandler *MinerHandler) MineBlock() {
	for {

		time.Sleep(minerHandler.genOpBlockTimeout)

		// 1 : last block
		lastblock := minerHandler.chainhandler.GetLongestValidChain()

		// 2 : new operation only
		//Todo : needs some lock in case of race condition
		newOperations := minerHandler.operationHandler.GetNewOperations()

		switch len(newOperations) {
		case 0:
			minerHandler.generateNoOpBlock(lastblock)
		default:
			minerHandler.generateOpBlock(newOperations, lastblock)
		}
	}
}

func (minerHandler *MinerHandler) generateOpBlock(newOperations []*entity.Operation, lastblock *entity.Block) {

	newOpBlock := entity.NewOpBlock(lastblock, newOperations)

	minerHandler.sharedchannel.Block <- newOpBlock
}

func (minerHandler *MinerHandler) generateNoOpBlock(lastblock *entity.Block) {

	newNoOpBlock := entity.NewNoOpBlock(lastblock)

	minerHandler.sharedchannel.Block <- newNoOpBlock
}

// func (minerHandler *MinerHandler) AddNoOpBlock(ctx context.Context) {

// 	for {
// 		time.Sleep(time.Second)

// 		lastblock := minerHandler.chainhandler.GetLongestValidChain()

// 		newNoOpBlock := &entity.Block{PrevHash: secsuit.ComputeHash(lastblock.String())}

// 		select {
// 		case <-ctx.Done():
// 			return
// 		default:
// 			minerHandler.chainhandler.AddBlock(newNoOpBlock)
// 		}
// 	}
// }

// func (minerHandler *MinerHandler) GenerateNoOpBlock() {
// 	contxt, cancelFunc := context.WithCancel(context.TODO())
// 	minerHandler.AddNoOpBlock(contxt)

// 	for {
// 		msg := <-minerHandler.cancelNoOpBlockGen
// 		if msg == 1 {
// 			cancelFunc()
// 		}

// 		if msg == 2 {
// 			contxt, cancelFunc = context.WithCancel(context.TODO())
// 			minerHandler.AddNoOpBlock(contxt)
// 		}
// 	}
// }

// func (minerHandler *MinerHandler) AddNoOpBlock(ctx context.Context) {

// 	for {
// 		time.Sleep(time.Second)

// 		lastblock := minerHandler.chainhandler.GetLongestValidChain()

// 		newNoOpBlock := &entity.Block{PrevHash: secsuit.ComputeHash(lastblock.String())}

// 		select {
// 		case <-ctx.Done():
// 			return
// 		default:
// 			minerHandler.chainhandler.AddBlock(newNoOpBlock)
// 		}
// 	}
// }
