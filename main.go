package main

import (
	test "Tests"
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
	fmt.Println("vim-go")
	rightMostCh = CreateLayers()
	mydata := &MyData{user: "Raju", data: nil, timeout: 20, val: 10}
	HttpRequest(mydata)
	mydata = &MyData{user: "Ravi", data: nil, timeout: 10, val: 10}
	HttpRequest(mydata)
	mydata = &MyData{user: "varun", data: nil, timeout: 15, val: 10}
	HttpRequest(mydata)
	for {
	}
}

func HttpRequest(inputdata *MyData) string {
	rightMostCh <- inputdata
	return fmt.Sprintf("Data Submitted to Queue: Take ur seat, Your number is :%d", len(rightMostCh))
}

func CreateLayers() chan *MyData {

	inputChan := make(chan *MyData)
	outputChan := CreateContext(inputChan)
	outputChan = Audit(outputChan)
	outputChan = RateLimitRequests(outputChan)
	outputChan = Authentication(outputChan)
	outputChan = SwitchLayer(outputChan)

	go PerformFullTests(outputChan)
	// Send data on Right most i.e First channel input
	// Out from the last data
	return inputChan
}

func PerformFullTests(rChan chan *MyData) {
	for data := range rChan {
		if data.success {
			fmt.Println("Do Full Tests on", data)
		}
		// Check the Success Flag and send it to next
		// stage of full tests or to error handling
		// pipeline
	}
}
func ProcessUser(data *MyData, lchan chan *MyData) {
	fmt.Println("Process User", data)
	data.success = false
	// Get Free Switch IP
	// Login and do snapl-load-
	// Try to see whether it is up every tick (10 seconds)
	// Verify Basic Tests
	// After full  time out, still not up, then the image is BAD.
	newTicker := time.NewTicker(time.Second * 2)
	// This is FAN-IN Design Pattren. Input from any of the
	// two channel to single output channel
	defer func() {
		newTicker.Stop()
		lchan <- data
	}()

	ftimeout := time.After(time.Second * data.timeout)
	for {
		select {
		case <-ftimeout:
			fmt.Println("TImeout", data)
			return
		case <-newTicker.C:
			fmt.Println("Check the status of switch every tick time", data)
			// if status is UP, Stop Ticker TODO and then do basic Tests
			test.BasicTests()
			//If Basic Tests Pass, sucess = true, else post onto
			// global Error Handling Process. Exist this loop in sucess
			// or failure
		}
	}
}

func SwitchLayer(rChan chan *MyData) chan *MyData {

	lChan := make(chan *MyData)

	// Wait input after authentication - Workers
	go func(lChan chan *MyData, rChan chan *MyData) {
		for newval := range rChan {
			go ProcessUser(newval, lChan)
		}
	}(lChan, rChan)

	return lChan
}

func Audit(rChan chan *MyData) chan *MyData {
	lChan := make(chan *MyData)

	go func(lChan chan *MyData, rChan chan *MyData) {
		for val := range rChan {
			fmt.Println("Operating in Audit ")
			lChan <- val
		}
	}(lChan, rChan)
	return lChan
}
func Authentication(rChan chan *MyData) chan *MyData {
	lChan := make(chan *MyData)

	go func(lChan chan *MyData, rChan chan *MyData) {
		for val := range rChan {
			fmt.Println("Operating in Authentication layer")
			lChan <- val
		}
	}(lChan, rChan)
	return lChan
}
func RateLimitRequests(rChan chan *MyData) chan *MyData {
	lChan := make(chan *MyData)

	go func(lChan chan *MyData, rChan chan *MyData) {
		for val := range rChan {
			fmt.Println("Operating in Rate limiting Routine")
			lChan <- val
		}
	}(lChan, rChan)
	return lChan
}
func CreateContext(rChan chan *MyData) chan *MyData {
	lChan := make(chan *MyData)

	go func(lChan chan *MyData, rChan chan *MyData) {
		for val := range rChan {
			fmt.Println("Operationg in Create Context Routine")
			lChan <- val
		}
	}(lChan, rChan)
	return lChan
}
