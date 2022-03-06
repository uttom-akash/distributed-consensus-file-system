package main

import (
	"flag"
	"log"
	"net/http"
	"os"
	"os/signal"
	"rfs/config"
	"rfs/handler/minehandler"
	"rfs/handler/synchandler"
	"rfs/miner"
	"rfs/models/entity"
	"rfs/models/modelconst"
	"time"
)

type MinerHandler struct{}

func (handler *MinerHandler) ServeHTTP(rw http.ResponseWriter, req *http.Request) {

}

func NewMinerHandler() {

}

func main() {

	minerId := flag.Int("minerid", 2, "miner id")

	flag.Parse()

	config.NewSingletonConfigHandler(config.ConsoleArg{MinerId: *minerId})

	minerHttp := miner.NewMinerHttp()

	go minerHttp.ListenPeerMiners()

	go minerHttp.ConnectPeerMiners()

	minehandler := minehandler.NewSingletonMinerHandler()

	synchandler := synchandler.NewSingletonSyncHandler()

	synchandler.Sync()

	go func() {
		time.Sleep(time.Minute)
		minehandler.AddNewOperation(entity.NewOperation("first.txt", modelconst.CREATE_FILE, nil))
	}()

	go func() {
		time.Sleep(3 * time.Minute)
		minehandler.AddNewOperation(entity.NewOperation("first.txt", modelconst.APPEND_RECORD, []byte("Append please")))
	}()

	interruptChan := make(chan os.Signal, 1)
	signal.Notify(interruptChan, os.Interrupt)
	signal.Notify(interruptChan, os.Kill)

	sig := <-interruptChan

	log.Println("Clossing ", sig)

	// // a := time.Now()

	// minerServer := http.Server{
	// 	Addr:    ":8080",
	// 	Handler: &MinerHandler{},
	// }

	// go func() {
	// 	err := minerServer.ListenAndServe()
	// 	log.Println("Started the server on port: ", 8080)

	// 	if err != nil {
	// 		log.Fatalf("Error : ", err)
	// 	}
	// }()

	// interruptChan := make(chan os.Signal, 1)

	// signal.Notify(interruptChan, os.Interrupt)
	// signal.Notify(interruptChan, os.Kill)

	// sig := <-interruptChan
	// log.Println("Got Interrupt: ", sig)

	// ctx, _ := context.WithTimeout(context.TODO(), 30*time.Second)
	// minerServer.Shutdown(ctx)
}
