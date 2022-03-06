package minehandler

import (
	"fmt"
	"rfs/handler/chainhandler"
	"rfs/handler/operationhandler"
	"rfs/models/entity"
	"sync"
	"time"
)

type MinerHandler struct {
	genOpBlockTimeout  time.Duration
	genNoOpBlockTime   time.Time
	cancelNoOpBlockGen chan int

	operationHandler *operationhandler.OperationHandler
	chainhandler     *chainhandler.ChainHandler
}

func NewMinerHandler() *MinerHandler {

	minerHandler := &MinerHandler{
		genOpBlockTimeout:  2 * time.Minute,
		genNoOpBlockTime:   time.Now(),
		cancelNoOpBlockGen: make(chan int),

		operationHandler: operationhandler.NewSingletonOperationHandler(),
		chainhandler:     chainhandler.NewSingletonChainHandler(),
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
	minerHandler.operationHandler.OperationChan <- operation
}

func (minerHandler *MinerHandler) MineBlock() {
	for {

		time.Sleep(minerHandler.genOpBlockTimeout)

		//Todo : needs some lock in case of race condition
		newOperations := minerHandler.operationHandler.GetNewOperations()

		switch len(newOperations) {
		case 0:
			minerHandler.generateNoOpBlock()
		default:
			minerHandler.generateOpBlock(newOperations)
		}
	}
}

func (minerHandler *MinerHandler) generateOpBlock(newOperations []*entity.Operation) {

	lastblock := minerHandler.chainhandler.GetLongestValidChain()

	newOpBlock := entity.NewOpBlock(lastblock, newOperations)

	minerHandler.chainhandler.Addblockchan <- newOpBlock
}

func (minerHandler *MinerHandler) generateNoOpBlock() {

	lastblock := minerHandler.chainhandler.GetLongestValidChain()

	newNoOpBlock := entity.NewNoOpBlock(lastblock)

	minerHandler.chainhandler.Addblockchan <- newNoOpBlock
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
