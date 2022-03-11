package sharedchannel

import (
	"fmt"
	"rfs/models/entity"
	"sync"
)

type SharedChannel struct {
	BroadcastBlock     chan *entity.Block
	Block              chan *entity.Block
	BroadcastOperation chan *entity.Operation
	Operation          chan *entity.Operation
}

func NewSharedChannel() *SharedChannel {
	return &SharedChannel{
		BroadcastBlock:     make(chan *entity.Block, 1),
		Block:              make(chan *entity.Block, 1),
		BroadcastOperation: make(chan *entity.Operation, 1),
		Operation:          make(chan *entity.Operation, 1),
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
