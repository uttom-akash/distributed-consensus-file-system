package pow

import "cfs/models/entity"

type IProofOfWork interface {
	DoProofWork(block *entity.Block, minDifficultyLvl int) int
}
