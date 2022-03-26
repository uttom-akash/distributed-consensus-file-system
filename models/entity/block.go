package entity

import (
	"cfs/config"
	"cfs/secsuit"
	"strconv"
	"time"
)

type Block struct {
	PrevHash   string
	Operations []*Operation
	MinerID    int
	Nonce      int
	TimeStamp  int64
	SerialNo   int
}

func (block *Block) String() string {
	//Todo: Need to implement it appropriately like markle root
	str := ""
	str += " " + block.PrevHash
	str += " " + strconv.Itoa(block.MinerID)
	str += " " + strconv.Itoa(block.Nonce)
	str += " " + strconv.FormatInt(block.TimeStamp, 10)
	str += " " + strconv.Itoa(block.SerialNo)

	for _, operation := range block.Operations {
		str += " " + operation.String()
	}

	return str
}

func (block *Block) Hash() string {
	return secsuit.ComputeHash(block.String())
}

func NewOpBlock(prevblock *Block, operations []*Operation) *Block {

	config := config.GetSingletonConfigHandler()

	return &Block{
		PrevHash:   prevblock.Hash(),
		Operations: operations,
		MinerID:    config.MinerConfig.MinerId,
		TimeStamp:  time.Now().Unix(),
		SerialNo:   prevblock.SerialNo + 1,
	}
}

func NewNoOpBlock(prevblock *Block) *Block {

	config := config.GetSingletonConfigHandler()

	return &Block{
		PrevHash:  prevblock.Hash(),
		MinerID:   config.MinerConfig.MinerId,
		TimeStamp: time.Now().Unix(),
		SerialNo:  prevblock.SerialNo + 1,
	}
}

func CreateGenesisBlock() *Block {
	//Todo: Proper implement - like config
	return &Block{
		SerialNo: 1,
	}
}

func (block *Block) PowDifficulty() int {

	blockHash := block.Hash()
	hashLength := len(blockHash)
	difficulty := 0

	for ; difficulty < hashLength; difficulty++ {
		if blockHash[difficulty] != '0' {
			break
		}
	}

	return difficulty
}
