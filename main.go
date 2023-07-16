package main

import (
	"cfs/cfslib"
	"cfs/config"
	"cfs/corehandler/synchandler"
	"cfs/models/entity"
	"cfs/models/message"
	"cfs/models/modelconst"
	"cfs/sharedchannel"
	"flag"
	"log"
	"os"
	"strconv"
	"time"
)

func ConfigureLogger(minerId int) *os.File {

	// logfile, logFErr := os.OpenFile("./storage/logs/log"+strconv.Itoa(*minerId)+".log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	logfile, logFErr := os.Create("./storage/logs/log" + strconv.Itoa(minerId) + ".log")

	if logFErr != nil {
		log.Fatalf("Error opening log file: %v", logFErr)
	}

	log.SetOutput(logfile)
	log.SetFlags(log.Llongfile | log.Ldate | log.Ltime) // Todo: add it later | log.LUTC

	return logfile
}

func main() {

	minerId := flag.Int("minerid", 2, "miner id")
	flag.Parse()

	logfile := ConfigureLogger(*minerId)
	defer logfile.Close()

	config.NewSingletonConfigHandler(config.ConsoleArg{MinerId: *minerId})

	sharedchannel := sharedchannel.NewSingletonSharedChannel()

	//Todo: remove once testing is done
	go func() {
		time.Sleep(time.Duration(cfslib.Random(1, 10)) * time.Minute)
		sharedchannel.InternalOperationChannel <- message.NewOperationMsg(entity.NewOperation("first.txt", modelconst.CREATE_FILE, nil), message.ADD)
	}()
	//Todo: remove once testing is done
	go func() {
		time.Sleep(time.Duration(cfslib.Random(10, 20)) * time.Minute)
		sharedchannel.InternalOperationChannel <- message.NewOperationMsg(entity.NewOperation("first.txt", modelconst.APPEND_RECORD, []byte("Append please")), message.ADD)
	}()

	synchandler := synchandler.NewSingletonSyncHandler()

	synchandler.Sync()
}
