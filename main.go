package main

import (
	Am "AutomationManager"
	logger "Logging"
	"fmt"
	"time"
)

type MyData struct {
	user    string
	data    map[string]string
	timeout time.Duration // timeout in seconds
	val     int
	success bool
}

var rightMostCh chan *MyData

func main() {
	log := logger.Init()
	log.Debug("Inside Main")

	am := &Am.AutomationMgr{}
	am.InitPackages()

	rightMostCh = am.ingressObj.CreateIngressPipeline()
	mydata := &MyData{user: "Raju", data: nil, timeout: 20, val: 10}
	HttpRequest(mydata)
	mydata = &MyData{user: "Ravi", data: nil, timeout: 10, val: 10}
	HttpRequest(mydata)
	mydata = &MyData{user: "varun", data: nil, timeout: 15, val: 10}
	HttpRequest(mydata)
	for {
	}
}

func HttpRequest(inputdata interface{}) string {
	rightMostCh <- inputdata
	return fmt.Sprintf("Data Submitted to Queue: Take ur seat, Your number is :%d", len(rightMostCh))
}
