package User

import (
	"fmt"
	"time"
)

type User struct {
}

func Init() *User {

	fmt.Println("Per User")
	return &User{}
}
func (user *User) ProcessUser(data *MyData, lchan chan *MyData) {
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
