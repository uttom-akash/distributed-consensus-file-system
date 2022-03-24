package main

import (
	"flag"
	"rfs/bclib"
	"rfs/config"
	"rfs/handler/synchandler"
	"rfs/models/entity"
	"rfs/models/message"
	"rfs/models/modelconst"
	"rfs/sharedchannel"
	"time"
)

func main() {

	minerId := flag.Int("minerid", 2, "miner id")
	flag.Parse()
	config.NewSingletonConfigHandler(config.ConsoleArg{MinerId: *minerId})

	sharedchannel := sharedchannel.NewSingletonSharedChannel()

	go func() {
		time.Sleep(time.Duration(bclib.Random(1, 10)) * time.Minute)
		sharedchannel.Operation <- message.NewOperationMsg(entity.NewOperation("first.txt", modelconst.CREATE_FILE, nil), message.ADD)
	}()
	go func() {
		time.Sleep(time.Duration(bclib.Random(10, 20)) * time.Minute)
		sharedchannel.Operation <- message.NewOperationMsg(entity.NewOperation("first.txt", modelconst.APPEND_RECORD, []byte("Append please")), message.ADD)
	}()

	synchandler := synchandler.NewSingletonSyncHandler()

	synchandler.Sync()
}
