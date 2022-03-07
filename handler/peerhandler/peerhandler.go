package peerhandler

import (
	"encoding/json"
	"log"
	"net/http"
	"rfs/bclib"
	"rfs/config"
	"rfs/handler/chainhandler"
	"rfs/handler/minehandler"
	"rfs/models/entity"
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
	chain := handler.chainHandler.GetChain()
	encodedJson, err := json.Marshal(chain)

	if err != nil {
		log.Println("error: ", err)
	}

	log.Println("Encoded: ", encodedJson)

	rw.Write(encodedJson)
}

func (handler *PeerHandler) servePong(rw http.ResponseWriter, req *http.Request) {
	rw.Write([]byte("pong!"))
}

func (handler *PeerHandler) ListenOperation(rw http.ResponseWriter, req *http.Request) {
	operation := new(entity.Operation)
	json.NewDecoder(req.Body).Decode(operation)

	handler.minerHandler.AddNewOperation(operation)
}

func (handler *PeerHandler) ListenBlock(rw http.ResponseWriter, req *http.Request) {
	block := new(entity.Block)
	json.NewDecoder(req.Body).Decode(block)
	handler.chainHandler.Addblockchan <- block
}
