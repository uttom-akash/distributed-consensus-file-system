package entity

import "rfs/bclib"

type Chain map[string]*Block

type Tails []*Block

type ChildBlockHash map[string][]string

type BlockChain struct {
	Chain          Chain `json:"chain,omitempty"`
	Tails          Tails `json:"tails,omitempty"`
	ChildBlockHash ChildBlockHash
}

func NewBlockchain() *BlockChain {
	chain := make(Chain)
	tails := make(Tails, 0)
	childBlockHash := make(ChildBlockHash)

	genesisBlock := CreateGenesisBlock()

	chain[genesisBlock.Hash()] = genesisBlock
	tails = append(tails, genesisBlock)

	return &BlockChain{
		Chain:          chain,
		Tails:          tails,
		ChildBlockHash: childBlockHash,
	}
}

func (chain *BlockChain) AddBlock(block *Block) {
	chain.Chain[block.Hash()] = block
	chain.ChildBlockHash[block.PrevHash] = append(chain.ChildBlockHash[block.PrevHash], block.Hash())
}

// last valid block in longest chain
func (chain *BlockChain) LastValidBlock() *Block {
	var lastblock *Block

	queue := bclib.NewQueue()

	//Todo: Can be improved
	genesisBlock := CreateGenesisBlock()

	queue.Push(genesisBlock.Hash())

	for !queue.IsEmpty() {
		levelSize := queue.Size()

		for levelSize > 0 {

			currentBlockHash := queue.Front().(string)
			queue.Pop()

			for _, childBlock := range chain.ChildBlockHash[currentBlockHash] {
				queue.Push(childBlock)
			}

			lastblock = chain.Chain[currentBlockHash]

			levelSize--
		}
	}

	return lastblock
}
