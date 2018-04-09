package AutomationManager

import (
	"Rest"
	"Rpc"
	"Tests"
	"User"
	"fmt"
)

type AutomationMgr struct {
	log        *log.Log
	serverObj  *Server.server
	ingressObj *Ingress.Ingress
	egressObj  *Egress.Egress
	rpcObj     *Rpc.Rpc
	restObj    *Rest.Rest
	errObj     *Error.Err
	testObj    *Tests.Test
	userObj    *User.User
}

func (am *AutomationMgr) InitPackages() {
	am.serverObj = Server.Init(am)
	am.ingressObj = Ingress.Init(am)
	am.egressObj = Egress.Init(am)
	am.egressObj = User.Init(am)
	am.restObj = Rest.Init(am)
	am.rpcObj = Rpc.Init(am)
	am.testsObj = Tests.Init(am)
	am.errObj = Error.Init(am)
	am.userObj = User.Init(am)

}
func Init() {
	fmt.Println("Automation Manager Init")
	return &AutomationManager{}
}
