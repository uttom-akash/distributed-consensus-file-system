package chainhandler

import (
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

	for block := range chainhandler.sharedchannel.InternalBlockChan {

		log.Println("ChainHandler/AddBlock - Processing block", block)

		if !chainhandler.validateBlock(block) {

			log.Println("ChainHandler/AddBlock - invalid block ", block)

			continue
		}

		log.Println("ChainHandler/AddBlock - successfully validated block ", block)

		chainhandler.operationHandler.SetOperationsStatus(block.Operations, modelconst.PENDING)

		operationsTobeRemoved := chainhandler.GetOperationsTobeRemoved(block)
		chainhandler.operationHandler.RemoveOperations(operationsTobeRemoved)

		// operationsTobeConfirmed := chainhandler.GetOperationsTobeConfirmed(block) //Todo: notify client

		// chainhandler.operationHandler.SetOperationsStatus(operationsTobeConfirmed, modelconst.CONFIRMED)

		chainhandler.chain.AddBlock(block)

		chainhandler.sharedchannel.BroadcastBlockChan <- block

		log.Println("ChainHandler/AddBlock - successfully added block ", block)
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

	queue := bclib.NewQueue()

	//Todo: Can be improved
	genesisBlock := chainhandler.chain.GenesisBlock

	queue.Push(genesisBlock.Hash())

	for !queue.IsEmpty() {
		currentBlockHash := queue.Front().(string)
		queue.Pop()

		chainhandler.sharedchannel.InternalBlockChan <- pChain.BlockHashMapper[currentBlockHash]

		for _, childBlock := range pChain.BlockTree[currentBlockHash] {
			queue.Push(childBlock)
		}
	}

	log.Println("ChainHandler/MargeChain - Out ")
}

func (chainhandler *ChainHandler) GetOperationsTobeConfirmed(block *entity.Block) []*entity.Operation {

	log.Println("ChainHandler/GetOperationsTobeConfirmed - In ")

	config := config.GetSingletonConfigHandler()

	var operations []*entity.Operation

	iterator := block
	hasParent := true

	for i := 0; hasParent && i <= int(math.Max(float64(config.SettingsConfig.ConfirmsPerFileCreate), float64(config.SettingsConfig.ConfirmsPerFileAppend))); i++ {

		if i == int(config.SettingsConfig.ConfirmsPerFileAppend) {
			for _, o := range iterator.Operations {
				if o.OperationType == modelconst.APPEND_RECORD {
					operations = append(operations, o)
				}
			}
		}

		if i == int(config.SettingsConfig.ConfirmsPerFileAppend) {
			for _, o := range iterator.Operations {
				if o.OperationType == modelconst.CREATE_FILE {
					operations = append(operations, o)
				}
			}
		}

		iterator, hasParent = chainhandler.chain.BlockHashMapper[iterator.PrevHash]
	}

	log.Println("ChainHandler/GetOperationsTobeConfirmed - Out ")

	return operations
}

func (chainhandler *ChainHandler) GetOperationsTobeRemoved(block *entity.Block) []*entity.Operation {

	log.Println("ChainHandler/GetOperationsTobeRemoved - In ")

	config := config.GetSingletonConfigHandler()
	numberOfblockToCheck := int(math.Max(float64(config.SettingsConfig.ConfirmsPerFileCreate), float64(config.SettingsConfig.ConfirmsPerFileAppend)))

	var operations []*entity.Operation

	blockInLongestChain := chainhandler.chain.LastValidBlock()
	iterator := block
	hasParent := true

	if blockInLongestChain.Hash() == block.PrevHash {
		for i := 0; i <= numberOfblockToCheck && hasParent; i++ {
			operations = append(operations, iterator.Operations...)

			iterator, hasParent = chainhandler.chain.BlockHashMapper[iterator.PrevHash]
		}
	}

	log.Println("ChainHandler/GetOperationsTobeRemoved - Out ")

	return operations
}
