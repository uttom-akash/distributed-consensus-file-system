package chainhandler

import (
	"fmt"
	"math"
	"rfs/models/entity"
	"rfs/secsuit"
	"sync"
)

type ChainHandler struct {
	chain        *entity.BlockChain
	Addblockchan chan *entity.Block
}

func NewChainHandler() *ChainHandler {

	return &ChainHandler{
		chain:        entity.NewBlockchain(),
		Addblockchan: make(chan *entity.Block, 2),
	}
}

var lock = &sync.Mutex{}
var chainhandlerInstance *ChainHandler

func NewSingletonChainHandler() *ChainHandler {

	if chainhandlerInstance == nil {
		lock.Lock()
		defer lock.Unlock()

		if chainhandlerInstance == nil {
			fmt.Println("Creating single instance now.")
			chainhandlerInstance = NewChainHandler()
		} else {
			fmt.Println("Single instance already created.")
		}
	} else {
		fmt.Println("Single instance already created.")
	}

	return chainhandlerInstance
}

func (chainhandler *ChainHandler) GetLongestValidChain() *entity.Block {
	largest := math.MinInt64
	var block *entity.Block

	for _, blck := range chainhandler.chain.Tails {
		if largest < blck.SerialNo {
			largest = blck.SerialNo
			block = blck
		}
	}

	go chainhandler.AddBlock()

	return block
}

func (chainhandler *ChainHandler) AddBlock() error {

	for block := range chainhandler.Addblockchan {

		if !chainhandler.ValidateBlock(block) {
			continue
		}

		for index, tail := range chainhandler.chain.Tails {
			if secsuit.ComputeHash(tail.String()) == block.PrevHash {

				block.SerialNo = tail.SerialNo + 1
				chainhandler.chain.Tails[index] = block
				break
			}
		}
		chainhandler.chain.Chain[secsuit.ComputeHash(block.String())] = block
	}

	return nil
}

func (chainhandler *ChainHandler) ValidateBlock(block *entity.Block) bool {

	//Check that the nonce for the block is valid: PoW is correct and has the right difficulty.
	//Check that the previous block hash points to a legal, previously generated, block.

	if _, ok := chainhandler.chain.Chain[block.PrevHash]; ok {
		return true
	}

	return false
}

func (chainhandler *ChainHandler) GetChain() *entity.BlockChain {

	return chainhandler.chain
}
