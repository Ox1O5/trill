package utils

import (
	"encoding/json"
	"io/ioutil"
)

type GlobalObj struct {
	//TcpServer  trill.IServer
	Host             string
	TcpPort          int
	Name             string
	Version          string
	MaxPacketSize    uint32
	MaxConnection    uint32
	WorkerPoolSize   uint32
	MaxWorkerTaskLen uint32
	MaxMsgChanLen    uint32
}

var GlobalObject *GlobalObj

func (g *GlobalObj) Load() {
	data, err := ioutil.ReadFile("configs/trill.json")
	if err != nil {
		panic(err)
	}
	err = json.Unmarshal(data, &GlobalObject)
	if err != nil {
		panic(err)
	}
}

func init() {
	GlobalObject = &GlobalObj{
		Name:             "Trill server",
		Version:          "V0.4",
		TcpPort:          9090,
		Host:             "0.0.0.0",
		MaxConnection:    10000,
		MaxPacketSize:    4096,
		WorkerPoolSize:   8,
		MaxWorkerTaskLen: 1024,
		MaxMsgChanLen:    16,
	}
	GlobalObject.Load()
}
