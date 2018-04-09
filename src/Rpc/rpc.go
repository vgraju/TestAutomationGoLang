package Rpc

import "fmt"

type Rpc struct {
}

func Init() *Rpc {
	fmt.Println("RPC Init")
	return &Rpc{}
}
func RpcHandler() {
	fmt.Println("RPC Handler")
}
