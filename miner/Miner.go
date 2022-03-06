package miner

import (
	"log"
	"net/http"
	"rfs/bclib"
	"rfs/config"
	"time"
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

type MinerHttp struct{}

//Connect
func (handler *MinerHttp) ListenPeerMiners() {

	log.Println("ListenPeerMiners")

	config := config.GetSingletonConfigHandler()

	rootHandler := http.NewServeMux()

	rootHandler.Handle("/ping", NewPeerHandler())

	bclib.HttpListen(http.Server{
		Addr:    config.MinerConfig.IpAddress + ":" + config.MinerConfig.Port,
		Handler: rootHandler,
	})
}

func (handler *MinerHttp) ListenClients() {

	bclib.HttpListen(http.Server{
		Addr:    ":8081",
		Handler: NewClientHandler(),
	})

}

func (handler *MinerHttp) ConnectPeerMiners() {
	con := config.GetSingletonConfigHandler()

	for {
		for _, peerId := range con.MinerConfig.Peers {
			log.Println("Connecting Peer : ", peerId)
			peerconfig := config.GetConfig(peerId)
			_, err := http.Get("http://" + peerconfig.IpAddress + ":" + peerconfig.Port + "/ping")

			if err != nil {
				log.Println("Error : pinging ", peerId, err)
			}

			log.Println("Ping success: ")

			time.Sleep(3 * time.Second)
		}
	}

}

func NewMinerHttp() *MinerHttp {
	return &MinerHttp{}
}

type PeerHandler struct{}

func (handler *PeerHandler) ServeHTTP(rw http.ResponseWriter, req *http.Request) {

}

func NewPeerHandler() *PeerHandler {
	return &PeerHandler{}
}

type ClientHandler struct{}

func (handler *ClientHandler) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	rw.Write([]byte("pong!"))
}

func NewClientHandler() *ClientHandler {
	return &ClientHandler{}
}
