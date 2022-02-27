package entity

import (
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
