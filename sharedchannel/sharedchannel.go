package sharedchannel

import (
	"cfs/models/entity"
	"cfs/models/message"
	"fmt"
	"sync"
)

type SharedChannel struct {
	BroadcastBlockChan     chan *entity.Block
	InternalBlockChan      chan *entity.Block
	BroadcastOperationChan chan *entity.Operation
	InternalOperationChan  chan *message.OperationChanMsg
	ConfirmedOperationChan chan *entity.Operation
}

func NewSharedChannel() *SharedChannel {
	return &SharedChannel{
		BroadcastBlockChan:     make(chan *entity.Block, 1),
		InternalBlockChan:      make(chan *entity.Block, 1),
		BroadcastOperationChan: make(chan *entity.Operation, 1),
		InternalOperationChan:  make(chan *message.OperationChanMsg, 1),
		ConfirmedOperationChan: make(chan *entity.Operation, 2), //Todo : make buffer size 1 later
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
