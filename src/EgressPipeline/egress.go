package EgressPipeline

import "fmt"

const (
	EGRESS_MAX_BUFF_SIZE = 100
)

type Egress struct {
	EgressChan chan interface{}
	am         *AutomationManager
}

func Init(am *AutomationManager) *Egress {
	// Buffer Channel
	eObj := &Egress{EgressChan: make(chan interface{}, EGRESS_MAX_BUFF_SIZE), am: am}
	return eObj
}
func (eObj *Egress) EgressProcess() {
	fmt.Println("vim-go")
}

func (eObj *Egress) GetEgressChan() chan interface{} {
	return eObj.EgressChan
}
