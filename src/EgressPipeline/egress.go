package main

import "fmt"

const (
	EGRESS_MAX_BUFF_SIZE = 100
)

type Egress struct {
	EgressChan chan interface{}
}

func Init() {
	// Buffer Channel
	eObj := &Egress{EgressChan: make(chan interface{}, EGRESS_MAX_BUFF_SIZE)}
}
func (eObj *Egress) EgressProcess() {
	fmt.Println("vim-go")
}

func (eObj *Egress) GetEgressChan() chan interface{} {
	return eObj.EgressChan
}
