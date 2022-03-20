package entity

import "rfs/bclib"

type BlockHashMapper map[string]*Block

type Tails []*Block

type ChildBlockHash map[string][]string

type BlockChain struct {
	BlockHashMapper BlockHashMapper `json:"block_hash_mapper,omitempty"`
	Tails           Tails           `json:"-"`
	BlockTree       ChildBlockHash  `json:"block_tree,omitempty"`
	GenesisBlock    *Block          `json:"genesis_block,omitempty"`
}

func NewBlockchain() *BlockChain {
	blockHashMapper := make(BlockHashMapper)
	tails := make(Tails, 0)
	childBlockHash := make(ChildBlockHash)
	genesisBlock := CreateGenesisBlock()

	blockHashMapper[genesisBlock.Hash()] = genesisBlock
	tails = append(tails, genesisBlock)

	return &BlockChain{
		BlockHashMapper: blockHashMapper,
		Tails:           tails,
		BlockTree:       childBlockHash,
		GenesisBlock:    genesisBlock,
	}
}

func (chain *BlockChain) AddBlock(block *Block) {
	chain.BlockHashMapper[block.Hash()] = block
	chain.BlockTree[block.PrevHash] = append(chain.BlockTree[block.PrevHash], block.Hash())
}

// last valid block in longest chain
func (chain *BlockChain) LastValidBlock() *Block {
	var lastblock *Block

	queue := bclib.NewQueue()

	genesisBlock := chain.GenesisBlock

	queue.Push(genesisBlock.Hash())

	for !queue.IsEmpty() {
		levelSize := queue.Size()
		var strongestBlock *Block

		for levelSize > 0 {

			currentBlockHash := queue.Front().(string)
			queue.Pop()

			for _, childBlock := range chain.BlockTree[currentBlockHash] {
				queue.Push(childBlock)
			}

			//Todo : check if strongestBlock null
			if strongestBlock.PowDifficulty() < chain.BlockHashMapper[currentBlockHash].PowDifficulty() {
				strongestBlock = chain.BlockHashMapper[currentBlockHash]
			}

			levelSize--
		}

		//Todo : check if strongestBlock null
		if strongestBlock != nil {
			lastblock = strongestBlock
		}
	}

	return lastblock
}
