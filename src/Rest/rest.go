package Rest

import "fmt"

type Rest struct {
	am *AutomationManager
}

func Init(am *AutomationManager) *Rest {
	fmt.Println("Rest Init")
	return &Rest{am: am}
}
func RestHandler() {
	fmt.Println("Rest Handler")
}
