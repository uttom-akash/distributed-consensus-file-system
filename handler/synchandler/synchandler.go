package synchandler

import (
	"fmt"
	"rfs/handler/chainhandler"
	"rfs/handler/minehandler"
	"rfs/handler/minernetworkoperationhandler"
	"rfs/handler/operationhandler"
	"rfs/handler/peerhandler"
	"sync"
)

type SyncHandler struct {
	mineHandler           *minehandler.MinerHandler
	chainhandler          *chainhandler.ChainHandler
	operationhandler      *operationhandler.OperationHandler
	peerhandler           *peerhandler.PeerHandler
	minerNetworkOperation *minernetworkoperationhandler.MinerNetworkOperationHandler
}

func (syncHandler *SyncHandler) Sync() {

	go syncHandler.peerhandler.ListenPeer()

	go syncHandler.minerNetworkOperation.DownloadChain()

	go syncHandler.minerNetworkOperation.DisseminateOperations()

	go syncHandler.minerNetworkOperation.DisseminateBlocks()

	go syncHandler.operationhandler.ListenOperationChannel()

	go syncHandler.chainhandler.AddBlock()

	go syncHandler.mineHandler.MineBlock()

}

func NewSyncHandler() *SyncHandler {
	return &SyncHandler{
		mineHandler:           minehandler.NewSingletonMinerHandler(),
		chainhandler:          chainhandler.NewSingletonChainHandler(),
		operationhandler:      operationhandler.NewSingletonOperationHandler(),
		peerhandler:           peerhandler.NewSingletonPeerHandler(),
		minerNetworkOperation: minernetworkoperationhandler.NewSingletonMinerNetworkOperationHandler(),
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
