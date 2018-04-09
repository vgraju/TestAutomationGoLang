package Tests

import "fmt"

type Test struct {
}

func BasicTests() {
	fmt.Println("Inside Test package basic tests")
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
