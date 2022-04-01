package chainhandler

import (
	"cfs/cfslib"
	"cfs/config"
	"cfs/corehandler/operationhandler"
	"cfs/models/entity"
	"cfs/models/modelconst"
	"cfs/sharedchannel"
	"fmt"
	"log"
	"math"
	"sync"
)

type ChainHandler struct {
	chain            *entity.BlockChain
	sharedchannel    *sharedchannel.SharedChannel
	operationHandler operationhandler.IOperationHandler
}

func NewChainHandler() IChainHandler {

	return &ChainHandler{
		chain:            entity.NewBlockchain(),
		sharedchannel:    sharedchannel.NewSingletonSharedChannel(),
		operationHandler: operationhandler.NewSingletonOperationHandler(),
	}
}

var lock = &sync.Mutex{}
var chainhandlerInstance IChainHandler

func NewSingletonChainHandler() IChainHandler {

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

	for block := range chainhandler.sharedchannel.InternalBlockChannel {

		log.Println("ChainHandler/AddBlock - Processing block", block)

		if !chainhandler.validateBlock(block) {

			log.Println("ChainHandler/AddBlock - invalid block ", block)

			continue
		}

		log.Println("ChainHandler/AddBlock - successfully validated block ", block)

		chainhandler.chain.AddBlock(block)

		chainhandler.sharedchannel.BroadcastBlockChannel <- block

		log.Println("ChainHandler/AddBlock - successfully added block ", block)

		operationsTobeRemoved := chainhandler.GetOperationsTobeRemoved(block)
		chainhandler.operationHandler.RemoveOperations(operationsTobeRemoved)

		log.Println("ChainHandler/AddBlock - successfully removed operations")

		chainhandler.PushConfirmedOperations(block) //Todo: notify client

		log.Println("ChainHandler/AddBlock - successfully push confirmed operations")
	}

	return nil
}

func (chainhandler *ChainHandler) validateBlock(block *entity.Block) bool {

	//Check that the nonce for the block is valid: PoW is correct and has the right difficulty.
	//Check that the previous block hash points to a legal, previously generated, block.

	if _, alreadyAdded := chainhandler.chain.BlockHashMapper[block.Hash()]; alreadyAdded {
		log.Println("ChainHandler/ValidateBlock - block is already added =", block.Hash())
		return false
	}

	if _, hasPerent := chainhandler.chain.BlockHashMapper[block.PrevHash]; !hasPerent {
		log.Println("ChainHandler/ValidateBlock - block doesn't have any parent block ", block.PrevHash)
		return false
	}

	return true
}

func (chainhandler *ChainHandler) GetChain() *entity.BlockChain {

	return chainhandler.chain
}

func (chainhandler *ChainHandler) MargeChain(pChain *entity.BlockChain) {

	log.Println("ChainHandler/MargeChain - In ")

	queue := cfslib.NewQueue()

	//Todo: Can be improved
	genesisBlock := chainhandler.chain.GenesisBlock

	queue.Push(genesisBlock.Hash())

	for !queue.IsEmpty() {
		currentBlockHash := queue.Front().(string)
		queue.Pop()

		chainhandler.sharedchannel.InternalBlockChannel <- pChain.BlockHashMapper[currentBlockHash]

		for _, childBlock := range pChain.BlockTree[currentBlockHash] {
			queue.Push(childBlock)
		}
	}

	log.Println("ChainHandler/MargeChain - Out ")
}

func (chainhandler *ChainHandler) PushConfirmedOperations(block *entity.Block) {

	log.Println("ChainHandler/GetOperationsTobeConfirmed - In ")

	config := config.GetSingletonConfigHandler()

	iterator := block
	hasParent := true

	for i := 0; hasParent && i <= int(math.Max(float64(config.SettingsConfig.ConfirmsPerFileCreate), float64(config.SettingsConfig.ConfirmsPerFileAppend))); i++ {

		if i == int(config.SettingsConfig.ConfirmsPerFileAppend) {
			for _, o := range iterator.Operations {
				if o.OperationType == modelconst.APPEND_RECORD && o.MinerID == config.ConsoleConfig.MinerId {
					chainhandler.sharedchannel.ConfirmedOperationChannel <- o
				}
			}
		}

		if i == int(config.SettingsConfig.ConfirmsPerFileAppend) {
			for _, o := range iterator.Operations {
				if o.OperationType == modelconst.CREATE_FILE && o.MinerID == config.ConsoleConfig.MinerId {
					chainhandler.sharedchannel.ConfirmedOperationChannel <- o
				}
			}
		}

		iterator, hasParent = chainhandler.chain.BlockHashMapper[iterator.PrevHash]
	}

	log.Println("ChainHandler/GetOperationsTobeConfirmed - Out ")
}

func (chainhandler *ChainHandler) GetOperationsTobeRemoved(block *entity.Block) []*entity.Operation {

	log.Println("ChainHandler/GetOperationsTobeRemoved - In ")

	config := config.GetSingletonConfigHandler()
	numberOfblockToCheck := int(math.Max(float64(config.SettingsConfig.ConfirmsPerFileCreate), float64(config.SettingsConfig.ConfirmsPerFileAppend)))

	var operationsTobeRemoved []*entity.Operation

	blockInLongestChain := chainhandler.chain.LastValidBlock()
	iterator := block
	hasParent := true

	if blockInLongestChain.Hash() == block.Hash() {
		for i := 0; i <= numberOfblockToCheck && hasParent; i++ {
			operationsTobeRemoved = append(operationsTobeRemoved, iterator.Operations...)

			iterator, hasParent = chainhandler.chain.BlockHashMapper[iterator.PrevHash]
		}
	}

	log.Println("ChainHandler/GetOperationsTobeRemoved - Out ")

	return operationsTobeRemoved
}
