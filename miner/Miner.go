package miner

import (
	"net/http"
	"rfs/bclib"
)

type MinerNetworkOperation interface {
	DownloadChain()

	DisseminateOperations()

	DisseminateBlocks()
}

type MinerWoker interface {
	GenerateOpBlock()

	GenerateNoOpBlock()

	ValidateBlock()

	ValidateOperation()
}

type MinerConnect interface {
	ListenPeerMiners()

	ListenClients()

	ConnectPeerMiners()
}

type Miner interface {
	MinerNetworkOperation
	MinerWoker
	MinerConnect
}

type MinerHandler struct{}

//Connect
func (handler *MinerHandler) ListenPeerMiners() {

	bclib.HttpListen(http.Server{
		Addr:    ":8080",
		Handler: NewPeerHandler(),
	})
}

func (handler *MinerHandler) ListenClients() {

	bclib.HttpListen(http.Server{
		Addr:    ":8081",
		Handler: NewClientHandler(),
	})

}

func (handler *MinerHandler) ConnectPeerMiners() {

}

type PeerHandler struct{}

func (handler *PeerHandler) ServeHTTP(rw http.ResponseWriter, req *http.Request) {

}

func NewPeerHandler() *PeerHandler {
	return &PeerHandler{}
}

type ClientHandler struct{}

func (handler *ClientHandler) ServeHTTP(rw http.ResponseWriter, req *http.Request) {

}

func NewClientHandler() *ClientHandler {
	return &ClientHandler{}
}
