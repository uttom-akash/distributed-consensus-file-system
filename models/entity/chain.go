package entity

import (
	"rfs/secsuit"
)

type Chain map[string]*Block

type Tails []*Block

type BlockChain struct {
	Chain Chain `json:"chain,omitempty"`
	Tails Tails `json:"tails,omitempty"`
}

func NewBlockchain() *BlockChain {
	chain := make(Chain)
	tails := make(Tails, 0)

	genesisBlock := CreateGenesisBlock()

	chain[secsuit.ComputeHash(genesisBlock.String())] = genesisBlock
	tails = append(tails, genesisBlock)

	return &BlockChain{
		Chain: chain,
		Tails: tails,
	}
}
