package entity

import (
	"cfs/config"
	"cfs/models/modelconst"
	"cfs/secsuit"
	"strconv"
	"time"
)

type Operation struct {
	OperationId    string //need to be distributed
	FileName       string
	OperationType  modelconst.OperationType
	Record         [512]byte
	MinerID        int
	TimeStamp      int64
	SponsoredCoins uint8
}

func NewOperation(fname string, operationType modelconst.OperationType, record []byte) *Operation {

	config := config.GetSingletonConfigHandler()
	minerId := config.MinerConfig.MinerId
	coins := uint8(0)

	if operationType == modelconst.CREATE_FILE {
		coins = config.SettingsConfig.NumCoinsPerFileCreate
	}

	var record512 [512]byte
	copy(record512[:], record)

	return &Operation{
		OperationId:    strconv.Itoa(minerId) + "-" + time.Now().String(),
		FileName:       fname,
		OperationType:  operationType,
		Record:         record512,
		MinerID:        minerId,
		TimeStamp:      time.Now().Unix(),
		SponsoredCoins: coins,
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
	str += " " + strconv.FormatInt(op.TimeStamp, 10)
	str += " " + strconv.Itoa(int(op.SponsoredCoins))

	return str
}

func (op *Operation) Hash() string {
	return secsuit.ComputeHash(op.String())
}
