package Server

import "fmt"

type server struct {
	am *AutomationManager
}

func (sobj *server) ProcessMessage() {
	fmt.Println("Insdie ProcessMessage of Server package")
}

func Init(am *AutomationManager) *server {
	fmt.Println("Initalizing the Server Package")
	return &server{am: am}
}
