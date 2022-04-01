package sharedchannel

import (
	"cfs/models/entity"
	"cfs/models/message"
	"fmt"
	"sync"
)

type SharedChannel struct {
	BroadcastBlockChannel       chan *entity.Block
	InternalBlockChannel        chan *entity.Block
	BroadcastOperationChannel   chan *entity.Operation
	InternalOperationChannel    chan *message.OperationChanMsg
	ConfirmedOperationChannel   chan *entity.Operation
	ClientOperationQueueChannel chan *entity.Operation
}

func NewSharedChannel() *SharedChannel {
	return &SharedChannel{
		BroadcastBlockChannel:       make(chan *entity.Block, 1),
		InternalBlockChannel:        make(chan *entity.Block, 1),
		BroadcastOperationChannel:   make(chan *entity.Operation, 1),
		InternalOperationChannel:    make(chan *message.OperationChanMsg, 1),
		ConfirmedOperationChannel:   make(chan *entity.Operation, 2), //Todo : make buffer size 1 later
		ClientOperationQueueChannel: make(chan *entity.Operation, 5),
	}
}

var lock = &sync.Mutex{}
var singletonInstance *SharedChannel

func NewSingletonSharedChannel() *SharedChannel {

	if singletonInstance == nil {
		lock.Lock()
		defer lock.Unlock()

		if singletonInstance == nil {
			fmt.Println("Creating single instance now.")
			singletonInstance = NewSharedChannel()
		} else {
			fmt.Println("Single instance already created.")
		}
	} else {
		fmt.Println("Single instance already created.")
	}

	return singletonInstance
}
