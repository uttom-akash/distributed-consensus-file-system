package config

import (
	"fmt"
	"log"
	"sync"
)

type MinerConfig struct {
	MinerId   int
	IpAddress string
	Port      string
	Peers     []int
}

type ConsoleArg struct {
	MinerId int
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
		Peers:     []int{4, 2},
	},
	{
		MinerId:   4,
		IpAddress: "127.0.0.1",
		Port:      "8083",
	},
}

type ConfigHandler struct {
	MinerConfig   MinerConfig
	ConsoleConfig ConsoleArg
}

var lock = &sync.Mutex{}
var confighandlerInstance *ConfigHandler

func NewConfigHandler(consoleArg ConsoleArg) *ConfigHandler {
	minerConfig := GetConfig(consoleArg.MinerId)

	return &ConfigHandler{
		MinerConfig:   minerConfig,
		ConsoleConfig: consoleArg,
	}
}

func NewSingletonConfigHandler(consoleArg ConsoleArg) *ConfigHandler {

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

func GetSingletonConfigHandler() *ConfigHandler {
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
