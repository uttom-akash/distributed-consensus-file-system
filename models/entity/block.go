package entity

import (
	"rfs/bclib"
	"rfs/config"
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

	time.Sleep(time.Duration(bclib.Random(40, 60)) * time.Second)

	config := config.GetSingletonConfigHandler()

	return &Block{
		PrevHash:   secsuit.ComputeHash(prevblock.String()),
		Operations: operations,
		MinerID:    config.MinerConfig.MinerId,
		TimeStamp:  time.Now(),
		SerialNo:   prevblock.SerialNo + 1,
	}
}

func NewNoOpBlock(prevblock *Block) *Block {

	time.Sleep(time.Duration(bclib.Random(20, 40)) * time.Second)

	config := config.GetSingletonConfigHandler()

	return &Block{
		PrevHash:  secsuit.ComputeHash(prevblock.String()),
		MinerID:   config.MinerConfig.MinerId,
		TimeStamp: time.Now(),
		SerialNo:  prevblock.SerialNo + 1,
	}
}

func CreateGenesisBlock() *Block {
	//Todo: Proper implement - like config
	return &Block{
		SerialNo: 1,
	}
}
