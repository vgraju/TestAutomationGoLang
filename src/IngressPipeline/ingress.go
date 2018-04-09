package IngressPipeline

import "fmt"

type Ingress struct {
	am *AutomationManager
}

func Init(am *AutomationManager) *Ingress {
	fmt.Println("Ingress Init")
	return &Ingress{am: am}
}
func IngressHandler() {
	fmt.Println("Ingress Handler")
}

func (i *Ingress) CreateIngressPipeline() chan interface{} {

	inputChan := make(chan interface{})
	outputChan := i.CreateContext(inputChan)
	outputChan = i.Audit(outputChan)
	outputChan = i.RateLimitRequests(outputChan)
	outputChan = i.Authentication(outputChan)
	outputChan = i.SwitchLayer(outputChan)

	go i.am.userObj.PerformFullTests(outputChan)
	// Send data on Right most i.e First channel input
	// Out from the last data
	return inputChan
}

func (ingress *Ingress) SwitchLayer(rChan chan interface{}) chan interface{} {

	lChan := make(chan interface{})

	// Wait input after authentication - Workers
	go func(lChan chan interface{}, rChan chan interface{}) {
		for newval := range rChan {
			go ProcessUser(newval, lChan)
		}
	}(lChan, rChan)

	return lChan
}

func (ingress *Ingress) Audit(rChan chan interface{}) chan interface{} {
	lChan := make(chan interface{})

	go func(lChan chan interface{}, rChan chan interface{}) {
		for val := range rChan {
			fmt.Println("Operating in Audit ")
			lChan <- val
		}
	}(lChan, rChan)
	return lChan
}
func (ingress *Ingress) Authentication(rChan chan interface{}) chan interface{} {
	lChan := make(chan interface{})

	go func(lChan chan interface{}, rChan chan interface{}) {
		for val := range rChan {
			fmt.Println("Operating in Authentication layer")
			lChan <- val
		}
	}(lChan, rChan)
	return lChan
}
func (ingress *Ingress) RateLimitRequests(rChan chan interface{}) chan interface{} {
	lChan := make(chan interface{})

	go func(lChan chan interface{}, rChan chan interface{}) {
		for val := range rChan {
			fmt.Println("Operating in Rate limiting Routine")
			lChan <- val
		}
	}(lChan, rChan)
	return lChan
}
func (ingress *Ingress) CreateContext(rChan chan interface{}) chan interface{} {
	lChan := make(chan interface{})

	go func(lChan chan interface{}, rChan chan interface{}) {
		for val := range rChan {
			fmt.Println("Operationg in Create Context Routine")
			lChan <- val
		}
	}(lChan, rChan)
	return lChan
}
