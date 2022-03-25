package peerserver

import (
	"cfs/cfslib"
	"cfs/config"
	"cfs/corehandler/chainhandler"
	"cfs/models/entity"
	"cfs/models/message"
	"cfs/sharedchannel"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"sync"
)

type PeerServerHandler struct {
	chainHandler  chainhandler.IChainHandler
	sharedchannel *sharedchannel.SharedChannel
}

func NewPeerServerHandler() *PeerServerHandler {
	return &PeerServerHandler{
		chainHandler:  chainhandler.NewSingletonChainHandler(),
		sharedchannel: sharedchannel.NewSingletonSharedChannel(),
	}
}

var lock = &sync.Mutex{}
var singletonInstance IPeerServerHandler

func NewSingletonPeerServerHandler() IPeerServerHandler {

	if singletonInstance == nil {
		lock.Lock()
		defer lock.Unlock()

		if singletonInstance == nil {
			fmt.Println("Creating single instance now.")
			singletonInstance = NewPeerServerHandler()
		} else {
			fmt.Println("Single instance already created.")
		}
	} else {
		fmt.Println("Single instance already created.")
	}

	return singletonInstance
}

func (handler *PeerServerHandler) ListenPeer() {
	log.Println("PeerServerHandler/ListenPeer - setting up")

	config := config.GetSingletonConfigHandler()

	rootHandler := http.NewServeMux()

	rootHandler.HandleFunc("/ping", handler.servePong)
	rootHandler.HandleFunc("/downloadchain", handler.serveChainDownload)
	rootHandler.HandleFunc("/operation", handler.listenOperation)
	rootHandler.HandleFunc("/block", handler.listenBlock)

	cfslib.HttpListen(http.Server{
		Addr:    config.MinerConfig.IpAddress + ":" + config.MinerConfig.Port,
		Handler: rootHandler,
	})

	log.Println("PeerServerHandler/ListenPeer - setting done.")
}

func (handler *PeerServerHandler) serveChainDownload(rw http.ResponseWriter, req *http.Request) {
	log.Println("PeerServerHandler/serveChainDownload - in")

	chain := handler.chainHandler.GetChain()
	encodedJson, encodedErr := json.Marshal(chain)

	if encodedErr != nil {
		log.Fatalf("PeerServerHandler/serveChainDownload - error encoding chain: %s", encodedErr)
	}

	rw.Header().Set("Content-Type", "application/json")
	rw.WriteHeader(http.StatusOK)
	rw.Write(encodedJson)

	log.Println("PeerServerHandler/serveChainDownload - out")
}

func (handler *PeerServerHandler) servePong(rw http.ResponseWriter, req *http.Request) {
	log.Println("PeerServerHandler/servePong - in")

	rw.Write([]byte("pong!"))

}

func (handler *PeerServerHandler) listenOperation(rw http.ResponseWriter, req *http.Request) {
	log.Println("PeerServerHandler/ListenOperation - in")

	operation := new(entity.Operation)

	decodedErr := json.NewDecoder(req.Body).Decode(operation)

	if decodedErr != nil {
		log.Fatalf("PeerServerHandler/ListenOperation - error decoding operation: %s", decodedErr)
	}

	handler.sharedchannel.InternalOperationChan <- message.NewOperationMsg(operation, message.ADD)

	log.Println("PeerServerHandler/ListenOperation - operation is added to channel $handler.minerHandler.AddNewOperation$")

	rw.WriteHeader(http.StatusOK)

}

func (handler *PeerServerHandler) listenBlock(rw http.ResponseWriter, req *http.Request) {

	log.Println("PeerServerHandler/ListenBlock - in")

	block := new(entity.Block)

	decodedErr := json.NewDecoder(req.Body).Decode(block)

	if decodedErr != nil {
		log.Fatalf("PeerServerHandler/ListenBlock - error decoding block: %s", decodedErr)
	}

	handler.sharedchannel.InternalBlockChan <- block

	log.Println("PeerServerHandler/ListenBlock - block is added to channel $handler.chainHandler.Addblockchan$")

	rw.WriteHeader(http.StatusOK)
}
