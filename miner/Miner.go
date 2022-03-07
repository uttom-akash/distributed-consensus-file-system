package miner

import (
	"encoding/json"
	"log"
	"net/http"
	"rfs/bclib"
	"rfs/config"
	"rfs/handler/chainhandler"
	"rfs/models/entity"
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
	rootHandler.Handle("/downloadchain", NewHandleChainDownload())

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

func (handler *MinerHttp) DownloadChain() {

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

type HandleChainDownload struct {
	chainHandler *chainhandler.ChainHandler
}

func (handler *HandleChainDownload) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	chain := handler.chainHandler.GetChain()
	encodedJson, err := json.Marshal(chain)

	if err != nil {
		log.Println("error: ", err)
	}

	log.Println("Encoded: ", encodedJson)

	rw.Write(encodedJson)
}

func NewHandleChainDownload() *HandleChainDownload {
	return &HandleChainDownload{chainHandler: chainhandler.NewSingletonChainHandler()}
}

type ClientHandler struct{}

func (handler *ClientHandler) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	rw.Write([]byte("pong!"))
}

func NewClientHandler() *ClientHandler {
	return &ClientHandler{}
}
