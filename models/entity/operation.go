package entity

import (
	"rfs/config"
	"rfs/models/modelconst"
	"rfs/secsuit"
	"strconv"
	"time"
)

type Operation struct {
	OperationId   string //need to be distributed
	FileName      string
	OperationType modelconst.OperationType
	Record        [512]byte
	MinerID       int
	TimeStamp     time.Time
	State         modelconst.OperationState
}

func NewOperation(fname string, operationType modelconst.OperationType, record []byte) *Operation {

	config := config.GetSingletonConfigHandler()
	minerId := config.MinerConfig.MinerId
	var record512 [512]byte
	copy(record512[:], record)

	return &Operation{
		OperationId:   strconv.Itoa(minerId) + "-" + time.Now().String(),
		FileName:      fname,
		OperationType: operationType,
		Record:        record512,
		MinerID:       minerId,
		TimeStamp:     time.Now(),
		State:         modelconst.NEW,
	}
}

func (op *Operation) String() string {
	//Todo: proper implementation
	str := ""
	str += op.OperationId
	str += " " + op.FileName
	str += " " + op.OperationType.String()
	str += " " + string(op.Record[:])
	str += " " + strconv.Itoa(op.MinerID)
	str += " " + op.TimeStamp.String()

	return str
}

func (op *Operation) Hash() string {
	return secsuit.ComputeHash(op.String())
}
