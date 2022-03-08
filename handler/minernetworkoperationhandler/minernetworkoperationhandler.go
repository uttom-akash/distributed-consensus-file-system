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

	for {
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
}

func (handler *MinerNetworkOperationHandler) DisseminateOperations() {
	con := config.GetSingletonConfigHandler()

	for operation := range handler.NewOperationsChan {
		for _, peerId := range con.MinerConfig.Peers {

			encodedOperation, _ := json.Marshal(operation)

			log.Println("Disseminating to ", peerId, "; operation : ", operation)

			peerconfig := config.GetConfig(peerId)
			resp, err := http.Post("http://"+peerconfig.IpAddress+":"+peerconfig.Port+"/operation", "application/json", bytes.NewBuffer(encodedOperation))

			if err != nil {
				log.Println("Error Disseminating ", peerId, err)
				continue
			}

			chain := new(entity.BlockChain)
			er := json.NewDecoder(resp.Body).Decode(chain)

			if er != nil {
				log.Println("Error Disseminating ", peerId, er)
				continue
			}

			log.Println("Disseminating success: ")

			time.Sleep(3 * time.Second)
		}
	}
}

func (handler *MinerNetworkOperationHandler) DisseminateBlocks() {
	con := config.GetSingletonConfigHandler()

	for block := range handler.NewBlocksChan {
		for _, peerId := range con.MinerConfig.Peers {

			encodedOperation, _ := json.Marshal(block)

			log.Println("Disseminating to ", peerId, "; block : ", block)

			peerconfig := config.GetConfig(peerId)
			resp, err := http.Post("http://"+peerconfig.IpAddress+":"+peerconfig.Port+"/block", "application/json", bytes.NewBuffer(encodedOperation))

			if err != nil {
				log.Println("Error disseminating block ", peerId, err)
				continue
			}

			chain := new(entity.BlockChain)
			er := json.NewDecoder(resp.Body).Decode(chain)

			if er != nil {
				log.Println("Error disseminating block ", peerId, er)
				continue
			}

			log.Println("Success - disseminating block  ")

			time.Sleep(3 * time.Second)
		}
	}
}
