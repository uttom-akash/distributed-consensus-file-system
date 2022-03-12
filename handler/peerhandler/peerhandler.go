package peerhandler

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"rfs/bclib"
	"rfs/config"
	"rfs/handler/chainhandler"
	"rfs/handler/minehandler"
	"rfs/models/entity"
	"rfs/sharedchannel"
	"sync"
)

type IPeerHandler interface {
}

type PeerHandler struct {
	chainHandler  *chainhandler.ChainHandler
	minerHandler  *minehandler.MinerHandler
	sharedchannel *sharedchannel.SharedChannel
}

func NewPeerHandler() *PeerHandler {
	return &PeerHandler{
		chainHandler:  chainhandler.NewSingletonChainHandler(),
		sharedchannel: sharedchannel.NewSingletonSharedChannel(),
		minerHandler:  minehandler.NewSingletonMinerHandler(),
	}
}

var lock = &sync.Mutex{}
var singletonInstance *PeerHandler

func NewSingletonPeerHandler() *PeerHandler {

	if singletonInstance == nil {
		lock.Lock()
		defer lock.Unlock()

		if singletonInstance == nil {
			fmt.Println("Creating single instance now.")
			singletonInstance = NewPeerHandler()
		} else {
			fmt.Println("Single instance already created.")
		}
	} else {
		fmt.Println("Single instance already created.")
	}

	return singletonInstance
}

func (handler *PeerHandler) ListenPeer() {
	log.Println("PeerHandler/ListenPeer - setting up")

	config := config.GetSingletonConfigHandler()

	rootHandler := http.NewServeMux()

	rootHandler.HandleFunc("/ping", handler.servePong)
	rootHandler.HandleFunc("/downloadchain", handler.serveChainDownload)
	rootHandler.HandleFunc("/operation", handler.ListenOperation)
	rootHandler.HandleFunc("/block", handler.ListenBlock)

	bclib.HttpListen(http.Server{
		Addr:    config.MinerConfig.IpAddress + ":" + config.MinerConfig.Port,
		Handler: rootHandler,
	})

	log.Println("PeerHandler/ListenPeer - setting done.")
}

func (handler *PeerHandler) serveChainDownload(rw http.ResponseWriter, req *http.Request) {
	log.Println("PeerHandler/serveChainDownload - in")

	chain := handler.chainHandler.GetChain()
	encodedJson, encodedErr := json.Marshal(chain)

	if encodedErr != nil {
		log.Fatalf("PeerHandler/serveChainDownload - error encoding chain: %s", encodedErr)
	}

	rw.Header().Set("Content-Type", "application/json")
	rw.WriteHeader(http.StatusOK)
	rw.Write(encodedJson)

	log.Println("PeerHandler/serveChainDownload - out")
}

func (handler *PeerHandler) servePong(rw http.ResponseWriter, req *http.Request) {
	log.Println("PeerHandler/servePong - in")

	rw.Write([]byte("pong!"))

}

func (handler *PeerHandler) ListenOperation(rw http.ResponseWriter, req *http.Request) {
	log.Println("PeerHandler/ListenOperation - in")

	operation := new(entity.Operation)

	decodedErr := json.NewDecoder(req.Body).Decode(operation)

	if decodedErr != nil {
		log.Fatalf("PeerHandler/ListenOperation - error decoding operation: %s", decodedErr)
	}

	handler.minerHandler.AddNewOperation(operation)

	log.Println("PeerHandler/ListenOperation - operation is added to channel $handler.minerHandler.AddNewOperation$")

	rw.WriteHeader(http.StatusOK)

}

func (handler *PeerHandler) ListenBlock(rw http.ResponseWriter, req *http.Request) {

	log.Println("PeerHandler/ListenBlock - in")

	block := new(entity.Block)

	decodedErr := json.NewDecoder(req.Body).Decode(block)

	if decodedErr != nil {
		log.Fatalf("PeerHandler/ListenBlock - error decoding block: %s", decodedErr)
	}

	handler.sharedchannel.Block <- block

	log.Println("PeerHandler/ListenBlock - block is added to channel $handler.chainHandler.Addblockchan$")

	rw.WriteHeader(http.StatusOK)
}
