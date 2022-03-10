package minernetworkoperationhandler

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"rfs/config"
	"rfs/models/entity"
	"sync"
	"time"
)

type MinerNetworkOperation interface {
	DownloadChain()

	DisseminateOperations()

	DisseminateBlocks()
}

type MinerNetworkOperationHandler struct {
	NewOperationsChan chan *entity.Operation
	NewBlocksChan     chan *entity.Block
}

func NewMinerNetworkOperationHandler() *MinerNetworkOperationHandler {
	return &MinerNetworkOperationHandler{
		NewOperationsChan: make(chan *entity.Operation, 1),
		NewBlocksChan:     make(chan *entity.Block, 1),
	}
}

var lock = &sync.Mutex{}
var singletonInstance *MinerNetworkOperationHandler

func NewSingletonMinerNetworkOperationHandler() *MinerNetworkOperationHandler {

	if singletonInstance == nil {
		lock.Lock()
		defer lock.Unlock()

		if singletonInstance == nil {
			fmt.Println("Creating single instance now.")
			singletonInstance = NewMinerNetworkOperationHandler()
		} else {
			fmt.Println("Single instance already created.")
		}
	} else {
		fmt.Println("Single instance already created.")
	}

	return singletonInstance
}

func (handler *MinerNetworkOperationHandler) DownloadChain() {

	con := config.GetSingletonConfigHandler()

	for _, peerId := range con.MinerConfig.Peers {
		log.Println("Downloading from peer : ", peerId)
		peerconfig := config.GetConfig(peerId)
		resp, err := http.Get("http://" + peerconfig.IpAddress + ":" + peerconfig.Port + "/downloadchain")

		if err != nil {
			log.Println("Error : Downloading from peer : ", peerId, err)
			continue
		}

		chain := new(entity.BlockChain)
		er := json.NewDecoder(resp.Body).Decode(chain)

		if er != nil {
			log.Println("Error : Downloading from peer : ", peerId, er)
			continue
		}

		log.Println("Success: Downloading from peer : ", peerId)

		time.Sleep(3 * time.Second)
	}

}

func (handler *MinerNetworkOperationHandler) DisseminateOperations() {
	log.Println("MinerNetworkOperation/DisseminateOperations - In")

	con := config.GetSingletonConfigHandler()

	for operation := range handler.NewOperationsChan {
		for _, peerId := range con.MinerConfig.Peers {

			log.Println("MinerNetworkOperation/DisseminateOperations - disseminating to ", peerId, "; operation : ", operation)

			encodedOperation, encoedeErr := json.Marshal(operation)

			if encoedeErr != nil {
				log.Println("MinerNetworkOperation/DisseminateOperations - error encoding operation", peerId, encoedeErr)
				continue
			}

			peerconfig := config.GetConfig(peerId)
			resp, httpErr := http.Post("http://"+peerconfig.IpAddress+":"+peerconfig.Port+"/operation", "application/json", bytes.NewBuffer(encodedOperation))

			if httpErr != nil {
				log.Println("MinerNetworkOperation/DisseminateOperations - error Disseminating operation", peerId, httpErr)
				continue
			}

			log.Println("MinerNetworkOperation/DisseminateOperations - dissemination success ", resp.StatusCode)

			time.Sleep(3 * time.Second)
		}
	}

	log.Println("MinerNetworkOperation/DisseminateOperations - Out")
}

func (handler *MinerNetworkOperationHandler) DisseminateBlocks() {

	log.Println("MinerNetworkOperation/DisseminateBlocks - In")

	con := config.GetSingletonConfigHandler()

	for block := range handler.NewBlocksChan {
		for _, peerId := range con.MinerConfig.Peers {

			log.Println("MinerNetworkOperation/DisseminateBlocks - disseminating to ", peerId, "; block : ", block)

			encodedOperation, encodedErr := json.Marshal(block)

			if encodedErr != nil {
				log.Println("MinerNetworkOperation/DisseminateBlocks - error encoding block ", peerId, encodedErr)
				continue
			}

			peerconfig := config.GetConfig(peerId)
			resp, httpErr := http.Post("http://"+peerconfig.IpAddress+":"+peerconfig.Port+"/block", "application/json", bytes.NewBuffer(encodedOperation))

			if httpErr != nil {
				log.Println("MinerNetworkOperation/DisseminateBlocks - Error disseminating block ", peerId, httpErr)
				continue
			}

			log.Println("MinerNetworkOperation/DisseminateBlocks - dissemination success ", resp.StatusCode)

			time.Sleep(3 * time.Second)
		}
	}

	log.Println("MinerNetworkOperation/DisseminateBlocks - Out")
}
