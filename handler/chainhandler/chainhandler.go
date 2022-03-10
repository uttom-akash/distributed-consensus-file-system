package chainhandler

import (
	"fmt"
	"log"
	"rfs/handler/minernetworkoperationhandler"
	"rfs/handler/operationhandler"
	"rfs/models/entity"
	"sync"
)

type ChainHandler struct {
	chain                        *entity.BlockChain
	Addblockchan                 chan *entity.Block
	minerNetworkOperationHandler *minernetworkoperationhandler.MinerNetworkOperationHandler
	operationHandler             *operationhandler.OperationHandler
}

func NewChainHandler() *ChainHandler {

	return &ChainHandler{
		chain:                        entity.NewBlockchain(),
		Addblockchan:                 make(chan *entity.Block, 2),
		minerNetworkOperationHandler: minernetworkoperationhandler.NewSingletonMinerNetworkOperationHandler(),
		operationHandler:             operationhandler.NewSingletonOperationHandler(),
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
	log.Println("ChainHandler/GetLongestValidChain - In")

	// largest := math.MinInt64
	// var block *entity.Block

	// for _, blck := range chainhandler.chain.Tails {
	// 	if largest < blck.SerialNo {
	// 		largest = blck.SerialNo
	// 		block = blck
	// 	}
	// }

	block := chainhandler.chain.LastValidBlock()

	log.Println("ChainHandler/GetLongestValidChain - Out ", block)

	return block
}

func (chainhandler *ChainHandler) AddBlock() error {
	log.Println("ChainHandler/AddBlock - In")

	for block := range chainhandler.Addblockchan {

		log.Println("ChainHandler/AddBlock - Processing block", block)

		if !chainhandler.ValidateBlock(block) {

			log.Println("ChainHandler/AddBlock - invalid block ", block)

			continue
		}

		log.Println("ChainHandler/AddBlock - successfully validated block ", block)

		// Todo : needs to define better. suppose how do we add two block in same previous block
		// for index, tail := range chainhandler.chain.Tails {
		// 	if tail.Hash() == block.PrevHash {

		// 		block.SerialNo = tail.SerialNo + 1
		// 		chainhandler.chain.Tails[index] = block
		// 		break
		// 	}
		// }

		chainhandler.operationHandler.SetOperationsPending(block.Operations)
		chainhandler.chain.AddBlock(block)

		chainhandler.minerNetworkOperationHandler.NewBlocksChan <- block

		log.Println("ChainHandler/AddBlock - successfully added block ", block)
	}

	return nil
}

func (chainhandler *ChainHandler) ValidateBlock(block *entity.Block) bool {

	//Check that the nonce for the block is valid: PoW is correct and has the right difficulty.
	//Check that the previous block hash points to a legal, previously generated, block.

	if _, alreadyAdded := chainhandler.chain.Chain[block.Hash()]; alreadyAdded {
		return false
	}

	if _, hasPerent := chainhandler.chain.Chain[block.PrevHash]; !hasPerent {
		return false
	}

	return true
}

func (chainhandler *ChainHandler) GetChain() *entity.BlockChain {

	return chainhandler.chain
}
