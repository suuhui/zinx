package utils

import (
	"encoding/json"
	"github.com/suuhui/zinx/v0.4/ziface"
	"io/ioutil"
)

type  GlobalObj struct {
	TcpServer ziface.IServer
	Host string `json:"host"`
	TcpPort int `json:"port"`
	Name string `json:"name"`
	Version string `json:"version"`
	MaxPacketSize uint32 `json:"maxPacketSize"`
	MaxConn int `json:"maxConn"`
}

var GlobalObject *GlobalObj

func init() {
	GlobalObject = &GlobalObj{
		TcpServer:     nil,
		Host:          "0.0.0.0",
		TcpPort:       8888,
		Name:          "zinx server",
		Version:       "V0.4",
		MaxPacketSize: 4096,
		MaxConn:       12000,
	}
	GlobalObject.Reload()
}

func (globalObj *GlobalObj) Reload() {
	data, err := ioutil.ReadFile("v0.4/conf/zinx.json")
	if err != nil {
		panic(err)
	}

	err = json.Unmarshal(data, &GlobalObject)
	if err != nil {
		panic(err)
	}
}
