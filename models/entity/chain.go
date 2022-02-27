package entity

import (
	"rfs/secsuit"
	"time"
)

type Chain map[string]*Block

type Tails []*Block

type BlockChain struct {
	Chain Chain
	Tails Tails
}

func NewBlockchain() *BlockChain {
	chain := make(Chain)
	tails := make(Tails, 0)

	genesisBlock := createGenesisBlock()

	chain[secsuit.ComputeHash(genesisBlock.String())] = genesisBlock
	tails = append(tails, genesisBlock)

	return &BlockChain{
		Chain: chain,
		Tails: tails,
	}
}

func createGenesisBlock() *Block {
	return &Block{
		TimeStamp: time.Now(),
		SerialNo:  1,
	}
}
