package pow

import (
	"cfs/models/entity"
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"log"
	"sync"
	"time"
)

type ProofOfWork struct {
}

func NewProofOfWorkHandler() IProofOfWork {

	return &ProofOfWork{}
}

var lock = &sync.Mutex{}
var proofOfWorkInstance IProofOfWork

func NewSingletonProofOfWorkHandler() IProofOfWork {

	if proofOfWorkInstance == nil {
		lock.Lock()
		defer lock.Unlock()

		if proofOfWorkInstance == nil {
			fmt.Println("Creating single instance now.")
			proofOfWorkInstance = NewProofOfWorkHandler()
		} else {
			fmt.Println("Single instance already created.")
		}
	} else {
		fmt.Println("Single instance already created.")
	}

	return proofOfWorkInstance
}

func (workproof ProofOfWork) DoProofWork(block *entity.Block, minDifficultyLvl int) int {
	log.Println("ProofOfWork/DoProofWork - Inside with minimum difficulty level: ", minDifficultyLvl)

	nonce := 0
	bfore := time.Now()

	for {
		block.Nonce = nonce
		difficultyLevel := block.PowDifficulty()

		//log.Println("ProofOfWork/DoProofWork- nonce: ", nonce, " -difficulty level: ", difficultyLevel, " -minimum difficulty level: ", minDifficultyLvl)

		if difficultyLevel >= minDifficultyLvl {
			break
		}
		nonce++
	}

	log.Println("ProofOfWork/DoProofWork -nonce: ", nonce, " -minimum difficulty level: ", minDifficultyLvl, "- hash: ", block.Hash())
	log.Println("ProofOfWork/DoProofWork -End with duration: ", time.Since(bfore))

	// if minDifficultyLvl == 5 {
	// 	os.Exit(1)
	// }

	return nonce
}

func ComputeHash(text string) string {
	h := md5.New()
	h.Write([]byte(text))
	str := hex.EncodeToString(h.Sum(nil))
	return str
}
