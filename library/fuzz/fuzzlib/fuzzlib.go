package fuzzlib

import (
	"fmt"
	"gowebfuzz/library/network/gwfhttp"
	"time"
)

type Module interface {
	Run(reqp *gwfhttp.RequestPacket) error
}

type LuaModule struct {
	Path string
}

func NewLuaModule(path string) LuaModule {
	lua := LuaModule{}
	lua.Path = path
	return lua
}

func (lm LuaModule) Run(reqp *gwfhttp.RequestPacket) error {
	fmt.Println("Lua Module is Running")
	return nil
}

type GolangModule struct {
	Path string
}

func NewGolangModule(path string) GolangModule{
	golang := GolangModule{}
	golang.Path = path
	return golang
}
func (glm GolangModule) Run(reqp *gwfhttp.RequestPacket) error {
	fmt.Printf("Golang Module is Running:%s,%s\n",reqp.GetRequestURL(),glm.Path)
	time.Sleep(3*time.Second)
	fmt.Printf("Golang Module is Done:%s,%s\n",reqp.GetRequestURL(),glm.Path)
	return nil
}

type JsonModule struct {
	Path string
}

func NewJsonModule(path string) JsonModule{
	json := JsonModule{}
	json.Path = path
	return json
}

func (jm JsonModule) Run(reqp *gwfhttp.RequestPacket) error {
	fmt.Println("Json Module is Running")
	return nil
}

type YamlModule struct {
	Path string
}

func NewYamlModule(path string) YamlModule {
	yaml := YamlModule{}
	yaml.Path = path
	return yaml
}

func (ym YamlModule) Run(reqp *gwfhttp.RequestPacket) error {
	fmt.Println("Yaml Module is Running")
	return nil
}
