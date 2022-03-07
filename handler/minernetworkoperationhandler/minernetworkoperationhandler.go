package minernetworkoperationhandler

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"rfs/config"
	"rfs/models/entity"
	"time"
)

type MinerNetworkOperation interface {
	DownloadChain()

	DisseminateOperations()

	DisseminateBlocks()
}

type MinerNetworkOperationHandler struct {
	NewOperations chan *entity.Operation
	NewBlocks     chan *entity.Block
}

func (handler *MinerNetworkOperationHandler) DownloadChain() {

	con := config.GetSingletonConfigHandler()

	for {
		for _, peerId := range con.MinerConfig.Peers {
			log.Println("Connecting Peer : ", peerId)
			peerconfig := config.GetConfig(peerId)
			resp, err := http.Get("http://" + peerconfig.IpAddress + ":" + peerconfig.Port + "/downloadchain")

			if err != nil {
				log.Println("Error : pinging ", peerId, err)
				continue
			}

			chain := new(entity.BlockChain)
			er := json.NewDecoder(resp.Body).Decode(chain)

			if er != nil {
				log.Println("Error : pinging ", peerId, er)
				continue
			}

			log.Println("Ping success: ", chain)

			time.Sleep(3 * time.Second)
		}
	}
}

func (handler *MinerNetworkOperationHandler) DisseminateOperations() {
	con := config.GetSingletonConfigHandler()

	for operation := range handler.NewOperations {
		for _, peerId := range con.MinerConfig.Peers {

			encodedOperation, _ := json.Marshal(operation)

			log.Println("Connecting Peer : ", peerId)
			peerconfig := config.GetConfig(peerId)
			resp, err := http.Post("http://"+peerconfig.IpAddress+":"+peerconfig.Port+"/operation", "application/json", bytes.NewBuffer(encodedOperation))

			if err != nil {
				log.Println("Error : pinging ", peerId, err)
				continue
			}

			chain := new(entity.BlockChain)
			er := json.NewDecoder(resp.Body).Decode(chain)

			if er != nil {
				log.Println("Error : pinging ", peerId, er)
				continue
			}

			log.Println("Ping success: ", chain)

			time.Sleep(3 * time.Second)
		}
	}
}

func (handler *MinerNetworkOperationHandler) DisseminateBlocks() {
	con := config.GetSingletonConfigHandler()

	for block := range handler.NewBlocks {
		for _, peerId := range con.MinerConfig.Peers {

			encodedOperation, _ := json.Marshal(block)

			log.Println("Connecting Peer : ", peerId)
			peerconfig := config.GetConfig(peerId)
			resp, err := http.Post("http://"+peerconfig.IpAddress+":"+peerconfig.Port+"/block", "application/json", bytes.NewBuffer(encodedOperation))

			if err != nil {
				log.Println("Error : pinging ", peerId, err)
				continue
			}

			chain := new(entity.BlockChain)
			er := json.NewDecoder(resp.Body).Decode(chain)

			if er != nil {
				log.Println("Error : pinging ", peerId, er)
				continue
			}

			log.Println("Ping success: ", chain)

			time.Sleep(3 * time.Second)
		}
	}
}
