package Rest

import (
	am "AutomationManager"
	"fmt"
)

type Rest struct {
	am *am.AutomationManager
}

func Init(am *am.AutomationManager) *Rest {
	fmt.Println("Rest Init")
	return &Rest{am: am}
}
func RestHandler() {
	fmt.Println("Rest Handler")
}
