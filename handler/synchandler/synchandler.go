package synchandler

import (
	"fmt"
	"rfs/handler/chainhandler"
	"rfs/handler/minehandler"
	"rfs/handler/operationhandler"
	"sync"
)

type SyncHandler struct {
	mineHandler      *minehandler.MinerHandler
	chainhandler     *chainhandler.ChainHandler
	operationhandler *operationhandler.OperationHandler
}

func (syncHandler *SyncHandler) Sync() {

	go syncHandler.operationhandler.ListenOperationChannel()

	go syncHandler.chainhandler.AddBlock()

	go syncHandler.mineHandler.MineBlock()

}

func NewSyncHandler() *SyncHandler {
	return &SyncHandler{
		mineHandler:      minehandler.NewSingletonMinerHandler(),
		chainhandler:     chainhandler.NewSingletonChainHandler(),
		operationhandler: operationhandler.NewSingletonOperationHandler(),
	}
}

var lock = &sync.Mutex{}
var singletonInstance *SyncHandler

func NewSingletonSyncHandler() *SyncHandler {

	if singletonInstance == nil {
		lock.Lock()
		defer lock.Unlock()

		if singletonInstance == nil {
			fmt.Println("Creating single instance now.")
			singletonInstance = NewSyncHandler()
		} else {
			fmt.Println("Single instance already created.")
		}
	} else {
		fmt.Println("Single instance already created.")
	}

	return singletonInstance
}
