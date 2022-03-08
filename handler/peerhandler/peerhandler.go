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
	"sync"
)

type IPeerHandler interface {
}

type PeerHandler struct {
	chainHandler *chainhandler.ChainHandler
	minerHandler *minehandler.MinerHandler
}

func NewPeerHandler() *PeerHandler {
	return &PeerHandler{
		chainHandler: chainhandler.NewSingletonChainHandler(),
		minerHandler: minehandler.NewSingletonMinerHandler(),
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
}

func (handler *PeerHandler) serveChainDownload(rw http.ResponseWriter, req *http.Request) {
	log.Println("Serving chain download")

	chain := handler.chainHandler.GetChain()
	encodedJson, err := json.Marshal(chain)

	if err != nil {
		log.Fatalln("Error: chain download", err)
	}

	rw.Write(encodedJson)

	log.Println("Success : chain download")
}

func (handler *PeerHandler) servePong(rw http.ResponseWriter, req *http.Request) {
	log.Println("pong response")

	rw.Write([]byte("pong!"))

}

func (handler *PeerHandler) ListenOperation(rw http.ResponseWriter, req *http.Request) {
	log.Println("ListenOperation: New Operation arrival")

	operation := new(entity.Operation)

	err := json.NewDecoder(req.Body).Decode(operation)

	if err != nil {
		log.Fatalln("ListenOperation: error decoding")
	}

	log.Println("ListenOperation: decoded operation ", operation)

	handler.minerHandler.AddNewOperation(operation)
}

func (handler *PeerHandler) ListenBlock(rw http.ResponseWriter, req *http.Request) {
	log.Println("ListenBlock: New block arrival")

	block := new(entity.Block)

	err := json.NewDecoder(req.Body).Decode(block)

	if err != nil {
		log.Fatalln("ListenBlock: error decoding")
	}

	log.Println("ListenBlock: decoded  block: ", block)

	handler.chainHandler.Addblockchan <- block
}
