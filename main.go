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
	"time"
)

func main() {

	minerId := flag.Int("minerid", 2, "miner id")
	flag.Parse()
	config.NewSingletonConfigHandler(config.ConsoleArg{MinerId: *minerId})

	sharedchannel := sharedchannel.NewSingletonSharedChannel()

	go func() {
		time.Sleep(time.Duration(cfslib.Random(1, 10)) * time.Minute)
		sharedchannel.InternalOperationChan <- message.NewOperationMsg(entity.NewOperation("first.txt", modelconst.CREATE_FILE, nil), message.ADD)
	}()
	go func() {
		time.Sleep(time.Duration(cfslib.Random(10, 20)) * time.Minute)
		sharedchannel.InternalOperationChan <- message.NewOperationMsg(entity.NewOperation("first.txt", modelconst.APPEND_RECORD, []byte("Append please")), message.ADD)
	}()

	synchandler := synchandler.NewSingletonSyncHandler()

	synchandler.Sync()
}
