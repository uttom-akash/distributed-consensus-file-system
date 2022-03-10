package main

import (
	"flag"
	"rfs/bclib"
	"rfs/config"
	"rfs/handler/minehandler"
	"rfs/handler/synchandler"
	"rfs/models/entity"
	"rfs/models/modelconst"
	"time"
)

func main() {

	minerId := flag.Int("minerid", 2, "miner id")
	flag.Parse()
	config.NewSingletonConfigHandler(config.ConsoleArg{MinerId: *minerId})

	minehandler := minehandler.NewSingletonMinerHandler()

	go func() {
		time.Sleep(time.Duration(bclib.Random(1, 2)) * time.Minute)
		minehandler.AddNewOperation(entity.NewOperation("first.txt", modelconst.CREATE_FILE, nil))
	}()
	go func() {
		time.Sleep(time.Duration(bclib.Random(2, 4)) * time.Minute)
		minehandler.AddNewOperation(entity.NewOperation("first.txt", modelconst.APPEND_RECORD, []byte("Append please")))
	}()

	synchandler := synchandler.NewSingletonSyncHandler()

	synchandler.Sync()

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
