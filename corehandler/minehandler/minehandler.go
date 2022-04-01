package minehandler

import (
	"cfs/config"
	"cfs/corehandler/chainhandler"
	"cfs/corehandler/operationhandler"
	"cfs/models/entity"
	"cfs/pow"
	"cfs/sharedchannel"
	"fmt"
	"sync"
	"time"
)

type MinerHandler struct {
	config           config.Configuration
	operationHandler operationhandler.IOperationHandler
	chainhandler     chainhandler.IChainHandler
	proofOfWork      pow.IProofOfWork
	sharedchannel    *sharedchannel.SharedChannel
}

func NewMinerHandler() IMinerHandler {

	minerHandler := &MinerHandler{
		config:           *config.GetSingletonConfigHandler(),
		sharedchannel:    sharedchannel.NewSingletonSharedChannel(),
		operationHandler: operationhandler.NewSingletonOperationHandler(),
		proofOfWork:      pow.NewSingletonProofOfWorkHandler(),
		chainhandler:     chainhandler.NewSingletonChainHandler(),
	}

	return minerHandler
}

var lock = &sync.Mutex{}
var singletonInstance IMinerHandler

func NewSingletonMinerHandler() IMinerHandler {

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

func (minerHandler *MinerHandler) MineBlock() {
	for {

		time.Sleep(time.Duration(minerHandler.config.SettingsConfig.GenOpBlockTimeout) * time.Minute)

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

	minerHandler.proofOfWork.DoProofWork(newOpBlock, int(minerHandler.config.SettingsConfig.PowPerOpBlock))

	minerHandler.sharedchannel.InternalBlockChannel <- newOpBlock
}

func (minerHandler *MinerHandler) generateNoOpBlock(lastblock *entity.Block) {

	newNoOpBlock := entity.NewNoOpBlock(lastblock)

	minerHandler.proofOfWork.DoProofWork(newNoOpBlock, int(minerHandler.config.SettingsConfig.PowPerNoOpBlock))

	minerHandler.sharedchannel.InternalBlockChannel <- newNoOpBlock
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
