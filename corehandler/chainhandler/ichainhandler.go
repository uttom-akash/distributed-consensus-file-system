package chainhandler

import "cfs/models/entity"

type IChainHandler interface {
	GetChain() *entity.BlockChain
	GetLongestValidChain() *entity.Block
	AddBlock() error
	MargeChain(pChain *entity.BlockChain)
	PushConfirmedOperations(block *entity.Block)
	GetOperationsTobeRemoved(block *entity.Block) []*entity.Operation
}
