package entity

import (
	"rfs/models/modelconst"
	"strconv"
	"time"
)

type Operation struct {
	OperationId   string //need to be distributed
	Fname         string
	OperationType modelconst.OperationType
	Record        [512]byte
	MinerID       int
	TimeStamp     time.Time
	State         modelconst.OperationState
}

func NewOperation(fname string, operationType modelconst.OperationType, record []byte) *Operation {

	minerId := 1 //Todo : read from conf
	var record512 [512]byte
	copy(record512[:], record)

	return &Operation{
		OperationId:   strconv.Itoa(minerId) + "-" + time.Now().String(),
		Fname:         fname,
		OperationType: operationType,
		Record:        record512,
		MinerID:       minerId,
		TimeStamp:     time.Now(),
		State:         modelconst.NEW,
	}
}

func (op *Operation) String() string {
	return ""
}
