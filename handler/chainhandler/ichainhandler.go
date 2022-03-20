package chainhandler

import "rfs/models/entity"

type IChainHandler interface {
	GetChain() *entity.BlockChain
	GetLongestValidChain() *entity.Block
	AddBlock() error
	MargeChain(pChain *entity.BlockChain)
	GetOperationsTobeConfirmed(block *entity.Block) []*entity.Operation
}
