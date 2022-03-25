package peerclient

import (
	"bytes"
	"cfs/config"
	"cfs/corehandler/chainhandler"
	"cfs/models/entity"
	"cfs/sharedchannel"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"
)

type PeerClientHandler struct {
	sharedchannel *sharedchannel.SharedChannel
	chainHandler  chainhandler.IChainHandler
}

func NewPeerClientHandler() IPeerClientHandler {
	return &PeerClientHandler{
		sharedchannel: sharedchannel.NewSingletonSharedChannel(),
		chainHandler:  chainhandler.NewSingletonChainHandler(),
	}
}

var lock = &sync.Mutex{}
var singletonInstance IPeerClientHandler

func NewSingletonPeerClientHandler() IPeerClientHandler {

	if singletonInstance == nil {
		lock.Lock()
		defer lock.Unlock()

		if singletonInstance == nil {
			fmt.Println("Creating single instance now.")
			singletonInstance = NewPeerClientHandler()
		} else {
			fmt.Println("Single instance already created.")
		}
	} else {
		fmt.Println("Single instance already created.")
	}

	return singletonInstance
}

func (handler *PeerClientHandler) DownloadChain() {

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

		handler.chainHandler.MargeChain(chain)

		log.Println("Success: Downloading from peer : ", peerId)

		time.Sleep(3 * time.Second)
	}

}

func (handler *PeerClientHandler) DisseminateOperations() {
	log.Println("PeerClientHandler/DisseminateOperations - In")

	con := config.GetSingletonConfigHandler()

	for operation := range handler.sharedchannel.BroadcastOperationChan {
		for _, peerId := range con.MinerConfig.Peers {

			log.Println("PeerClientHandler/DisseminateOperations - disseminating to ", peerId, "; operation : ", operation)

			encodedOperation, encoedeErr := json.Marshal(operation)

			if encoedeErr != nil {
				log.Println("PeerClientHandler/DisseminateOperations - error encoding operation", peerId, encoedeErr)
				continue
			}

			peerconfig := config.GetConfig(peerId)
			resp, httpErr := http.Post("http://"+peerconfig.IpAddress+":"+peerconfig.Port+"/operation", "application/json", bytes.NewBuffer(encodedOperation))

			if httpErr != nil {
				log.Println("PeerClientHandler/DisseminateOperations - error Disseminating operation", peerId, httpErr)
				continue
			}

			log.Println("PeerClientHandler/DisseminateOperations - dissemination success ", resp.StatusCode)

			time.Sleep(3 * time.Second)
		}
	}

	log.Println("PeerClientHandler/DisseminateOperations - Out")
}

func (handler *PeerClientHandler) DisseminateBlocks() {

	log.Println("PeerClientHandler/DisseminateBlocks - In")

	con := config.GetSingletonConfigHandler()

	for block := range handler.sharedchannel.BroadcastBlockChan {
		for _, peerId := range con.MinerConfig.Peers {

			log.Println("PeerClientHandler/DisseminateBlocks - disseminating to ", peerId, "; block : ", block)

			encodedOperation, encodedErr := json.Marshal(block)

			if encodedErr != nil {
				log.Println("PeerClientHandler/DisseminateBlocks - error encoding block ", peerId, encodedErr)
				continue
			}

			peerconfig := config.GetConfig(peerId)
			resp, httpErr := http.Post("http://"+peerconfig.IpAddress+":"+peerconfig.Port+"/block", "application/json", bytes.NewBuffer(encodedOperation))

			if httpErr != nil {
				log.Println("PeerClientHandler/DisseminateBlocks - Error disseminating block ", peerId, httpErr)
				continue
			}

			log.Println("PeerClientHandler/DisseminateBlocks - dissemination success ", resp.StatusCode)

			time.Sleep(3 * time.Second)
		}
	}

	log.Println("PeerClientHandler/DisseminateBlocks - Out")
}
