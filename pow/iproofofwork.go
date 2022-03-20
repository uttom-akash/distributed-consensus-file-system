package pow

import "rfs/models/entity"

type IProofOfWork interface {
	DoProofWork(block *entity.Block, minDifficultyLvl int) int
}
