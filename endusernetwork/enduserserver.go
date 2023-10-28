package endusernetwork

import (
	"cfs/cfslib"
	"cfs/config"
	"cfs/endusernetwork/dtos"
	"cfs/models/entity"
	"cfs/models/message"
	"cfs/models/modelconst"
	"cfs/sharedchannel"
	"encoding/json"
	"log"
	"net/http"
)

type EndUserServer struct {
	sharedChannel *sharedchannel.SharedChannel
}

func NewPeerServerHandler() *EndUserServer {

	endUserServer := &EndUserServer{
		sharedChannel: sharedchannel.NewSingletonSharedChannel(),
	}

	return endUserServer
}

func (endUserServer *EndUserServer) Serve() {
	log.Println("EndUserServer - setting up")

	config := config.GetSingletonConfigHandler()
	rootHandler := http.NewServeMux()

	rootHandler.HandleFunc("/operation/ping", endUserServer.servePong)
	rootHandler.HandleFunc("/operation/file/create", endUserServer.createFile)
	rootHandler.HandleFunc("/operation/record/append", endUserServer.appendRecord)

	cfslib.HttpListen(http.Server{
		Addr:    config.MinerConfig.IpAddress + ":" + config.MinerConfig.EndUserNetworkPort,
		Handler: rootHandler,
	})

	log.Println("PeerServerHandler/ListenPeer - setting done.")
}

// CreateFile CreateFile(filename string) (err error)
func (endUserServer *EndUserServer) createFile(rw http.ResponseWriter, req *http.Request) {

	log.Println("EndUserServer/CreateFile - in")

	filename := req.URL.Query().Get("filename")

	log.Println(filename)

	endUserServer.sharedChannel.InternalOperationChannel <- message.NewOperationMsg(
		entity.NewOperation(
			filename, modelconst.CREATE_FILE, nil),
		message.ADD)

	rw.WriteHeader(http.StatusOK)

	log.Println("EndUserServer/CreateFile - out")
}

// AppendRecord AppendRec(filename string, record *Record) (recordNum uint16, err error)
func (endUserServer *EndUserServer) appendRecord(rw http.ResponseWriter, req *http.Request) {

	log.Println("EndUserServer/AppendRecord - in")

	filename := req.URL.Query().Get("filename")

	log.Println(filename)

	appendRecordCommand := new(dtos.AppendRecordCommand)
	decodedErr := json.NewDecoder(req.Body).Decode(appendRecordCommand)

	if decodedErr != nil {
		log.Fatalf("PeerServerHandler/ListenOperation - error decoding operation: %s", decodedErr)
	}

	endUserServer.sharedChannel.InternalOperationChannel <- message.NewOperationMsg(
		entity.NewOperation(
			filename,
			modelconst.APPEND_RECORD, []byte(appendRecordCommand.Record)),
		message.ADD)

	rw.WriteHeader(http.StatusOK)

	log.Println("EndUserServer/AppendRecord - out")
}

func (endUserServer *EndUserServer) servePong(rw http.ResponseWriter, req *http.Request) {

	log.Println("EndUserServer/servePong - in")

	rw.WriteHeader(http.StatusOK)

	log.Println("EndUserServer/servePong - out")
}
