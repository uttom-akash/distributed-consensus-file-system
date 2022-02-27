package chainhandler

import (
	"math"
	"rfs/models/entity"
	"rfs/secsuit"
)

type ChainHandler struct {
	chain *entity.BlockChain
}

func NewChainHandler() *ChainHandler {

	return &ChainHandler{
		chain: entity.NewBlockchain(),
	}
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

	return block
}

func (chainhandler *ChainHandler) AddBlock(block *entity.Block) error {

	if !chainhandler.ValidateBlock(block) {
		return nil
	}

	for index, tail := range chainhandler.chain.Tails {
		if secsuit.ComputeHash(tail.String()) == block.PrevHash {

			block.SerialNo = tail.SerialNo + 1
			chainhandler.chain.Tails[index] = block
			break
		}
	}
	chainhandler.chain.Chain[secsuit.ComputeHash(block.String())] = block

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
