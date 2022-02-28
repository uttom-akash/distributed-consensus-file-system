package entity

import (
	"rfs/secsuit"
	"strconv"
	"time"
)

type Block struct {
	PrevHash   string
	Operations []*Operation
	MinerID    int
	Nonce      int
	TimeStamp  time.Time
	SerialNo   int
}

func (block *Block) String() string {
	//Todo: Need to implement it appropriately
	return block.PrevHash + block.TimeStamp.String() + strconv.Itoa(block.SerialNo)
}

func NewOpBlock(prevblock *Block, operations []*Operation) *Block {
	return &Block{
		PrevHash:   secsuit.ComputeHash(prevblock.String()),
		Operations: operations,
		TimeStamp:  time.Now(),
	}
}

func NewNoOpBlock(prevblock *Block) *Block {
	return &Block{
		PrevHash:  secsuit.ComputeHash(prevblock.String()),
		TimeStamp: time.Now(),
	}
}
