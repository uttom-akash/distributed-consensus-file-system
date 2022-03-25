package chainhandler

import "cfs/models/entity"

type IChainHandler interface {
	GetChain() *entity.BlockChain
	GetLongestValidChain() *entity.Block
	AddBlock() error
	MargeChain(pChain *entity.BlockChain)
	GetOperationsTobeConfirmed(block *entity.Block) []*entity.Operation
	GetOperationsTobeRemoved(block *entity.Block) []*entity.Operation
}
