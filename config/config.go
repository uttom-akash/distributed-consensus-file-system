package config

import (
	"fmt"
	"log"
	"sync"
)

type MinerConfig struct {
	MinerId             int
	IpAddress           string //Todo : IncomingMinersAddr ip:port
	Port                string
	Peers               []int
	OutgoingMinersIP    string
	IncomingClientsAddr string
}

type ConsoleArg struct {
	MinerId int
}

type SettingsConfig struct {
	MinedCoinsPerOpBlock   uint8  // The number of record coins mined for an op block
	MinedCoinsPerNoOpBlock uint8  // The number of record coins mined for a no-op block
	NumCoinsPerFileCreate  uint8  // The number of record coins charged for creating a file
	GenOpBlockTimeout      uint8  // Time in milliseconds, the minimum time between op block mining (see diagram above).
	GenesisBlockHash       string // The genesis (first) block MD5 hash for this blockchain
	PowPerOpBlock          uint8  // The op block difficulty (proof of work setting: number of zeroes)
	PowPerNoOpBlock        uint8  // The no-op block difficulty (proof of work setting: number of zeroes)
	ConfirmsPerFileCreate  uint8  // The number of confirmations for a create file operation (the number of blocks that must follow the block containing a create file operation along longest chain before the CreateFile call can return successfully)
	ConfirmsPerFileAppend  uint8  // The number of confirmations for an append operation
}

var configs = []MinerConfig{
	{
		MinerId:   1,
		IpAddress: "127.0.0.1",
		Port:      "8080",
		Peers:     []int{2, 3},
	},
	{
		MinerId:   2,
		IpAddress: "127.0.0.1",
		Port:      "8081",
		Peers:     []int{3, 4},
	},
	{
		MinerId:   3,
		IpAddress: "127.0.0.1",
		Port:      "8082",
		Peers:     []int{4, 1},
	},
	{
		MinerId:   4,
		IpAddress: "127.0.0.1",
		Port:      "8083",
		Peers:     []int{1},
	},
}

type Configuration struct {
	MinerConfig    MinerConfig
	ConsoleConfig  ConsoleArg
	SettingsConfig SettingsConfig
}

var lock = &sync.Mutex{}
var confighandlerInstance *Configuration

func NewConfigHandler(consoleArg ConsoleArg) *Configuration {
	return &Configuration{
		MinerConfig:    GetConfig(consoleArg.MinerId),
		ConsoleConfig:  consoleArg,
		SettingsConfig: GetSettingsConfig(),
	}
}

func NewSingletonConfigHandler(consoleArg ConsoleArg) *Configuration {

	if confighandlerInstance == nil {
		lock.Lock()
		defer lock.Unlock()

		if confighandlerInstance == nil {
			fmt.Println("Creating single instance now.")
			confighandlerInstance = NewConfigHandler(consoleArg)
		} else {
			fmt.Println("Single instance already created.")
		}
	} else {
		fmt.Println("Single instance already created.")
	}

	return confighandlerInstance
}

func GetSingletonConfigHandler() *Configuration {
	if confighandlerInstance == nil {
		log.Fatal("ConfigHandler is not created yet.")
	}

	return confighandlerInstance
}

func GetConfig(id int) MinerConfig {
	if id <= 0 {
		log.Println("miner Id ", id)
	}
	return configs[id-1]
}

func GetSettingsConfig() SettingsConfig {

	return SettingsConfig{
		MinedCoinsPerOpBlock:   3,
		MinedCoinsPerNoOpBlock: 3,
		NumCoinsPerFileCreate:  2,
		GenOpBlockTimeout:      5,
		GenesisBlockHash:       "",
		PowPerOpBlock:          5,
		PowPerNoOpBlock:        5,
		ConfirmsPerFileCreate:  3,
		ConfirmsPerFileAppend:  5,
	}
}
