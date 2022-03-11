package synchandler

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/signal"
	"rfs/config"
	"rfs/handler/chainhandler"
	"rfs/handler/minehandler"
	"rfs/handler/minernetworkoperationhandler"
	"rfs/handler/operationhandler"
	"rfs/handler/peerhandler"
	"strconv"
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

	syncHandler.shutdownOnInterrupt()
}

var lock = &sync.Mutex{}
var singletonInstance *SyncHandler

func NewSyncHandler() *SyncHandler {
	return &SyncHandler{
		mineHandler:           minehandler.NewSingletonMinerHandler(),
		chainhandler:          chainhandler.NewSingletonChainHandler(),
		operationhandler:      operationhandler.NewSingletonOperationHandler(),
		peerhandler:           peerhandler.NewSingletonPeerHandler(),
		minerNetworkOperation: minernetworkoperationhandler.NewSingletonMinerNetworkOperationHandler(),
	}
}

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

func (syncHandler *SyncHandler) shutdownOnInterrupt() {

	interruptChan := make(chan os.Signal, 1)
	signal.Notify(interruptChan, os.Interrupt)
	signal.Notify(interruptChan, os.Kill)

	sig := <-interruptChan

	syncHandler.writeChain()

	log.Println("Shutting down the system gracefully ", sig)
}

func (syncHandler *SyncHandler) writeChain() {

	log.Println("SyncHandler/write - Writing")

	config := config.GetSingletonConfigHandler()

	jsonFile, ioErr := os.Create("./Chain" + strconv.Itoa(config.MinerConfig.MinerId) + ".json")

	if ioErr != nil {
		log.Println("SyncHandler/write - error creating json file ", ioErr)
		return
	}

	defer jsonFile.Close()

	jsonData, encodedErr := json.Marshal(syncHandler.chainhandler.GetChain().BlockHashMapper)

	if encodedErr != nil {
		log.Println("SyncHandler/write - error encoding ", encodedErr)
		return
	}

	// sanity check
	fmt.Println("SyncHandler/write - ", string(jsonData))

	jsonFile.Write(jsonData)
	jsonFile.Close()

	fmt.Println("SyncHandler/write - JSON data written to ", jsonFile.Name())
}
