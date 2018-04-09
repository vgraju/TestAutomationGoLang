package main

import (
	"fmt"
	"io"
	"regexp"
	"strconv"
	"strings"
	"time"

	"golang.org/x/crypto/ssh"
)

const (
	CREATE_RECORD int = iota
	INGRESS_UPDATE_RECORD
	EGRESS_SUCCESS_RECORD
	EGRESS_RETRY_RECORD
	DELETE_RECORD
)

type Request struct {
	cmd   string
	mod   string
	rcvCh chan *Response
}
type Response struct {
	cmd    string
	output string
	mod    string
}
type MySsh struct {
	client          *ssh.Client
	r               io.Reader
	w               io.Writer
	session         *ssh.Session
	IngressInputCh  chan interface{}
	EgressInputCh   chan interface{}
	ReadCh          chan string
	RecycleDuration time.Duration
	MyMap           map[string]map[string]interface{}
}

func main() {

	elem := MySsh{RecycleDuration: 5}

	elem.MyMap = make(map[string]map[string]interface{})
	elem.ReadCh = make(chan string)
	elem.IngressInputCh = make(chan interface{}, 10)
	elem.EgressInputCh = make(chan interface{}, 10)

	config := &ssh.ClientConfig{User: "root", HostKeyCallback: ssh.InsecureIgnoreHostKey(), Timeout: 0}
	fmt.Println("Client:", config)

	client, err := ssh.Dial("tcp", "192.168.100.105:22", config)
	if err != nil {
		fmt.Println("Connedtion failed :", err)
	} else {

		fmt.Println("Connedtion Success:", err)
		elem.client = client
	}

	leftChan := elem.CreateIngressPipeline(elem.IngressInputCh)
	finalCh := elem.CreateEgressPipeline(leftChan)

	slice := []string{"intf", "asic"}

	for _, mod := range slice {
		str := "docker ps | grep " + mod + " | grep opt | wc"
		elem.SendCommand(mod, str)
	}

	finalGo(finalCh)

}

func finalGo(finalCh chan interface{}) {
	for val := range finalCh {
		if val == "end" {
			fmt.Println("End: ", val)
		}
		fmt.Println("Fina Reveive Value ", val)
	}
}
func (c *MySsh) CreateIngressPipeline(rightChan chan interface{}) chan interface{} {

	var leftChan chan interface{}

	leftChan = c.SaveContext(rightChan)
	leftChan = c.Audit(leftChan)
	leftChan = c.RecordIngress(leftChan)
	leftChan = c.Final(leftChan) // THis will put the message on teh caller queue
	return leftChan
}
func (c *MySsh) CreateEgressPipeline(rightChan chan interface{}) chan interface{} {

	var leftChan chan interface{}

	leftChan = c.RecordEgress(rightChan)
	leftChan = c.ProcessEgress(leftChan)
	leftChan = c.FinalEgress(leftChan)
	return leftChan
}

func (c *MySsh) RecordMap(cmd int, mod string) {
	switch cmd {
	case INGRESS_UPDATE_RECORD:
		// Use Locks here
		if _, ok := c.MyMap[mod]; !ok {
			fmt.Println("Creating Map element for mod :", mod)
			c.MyMap[mod] = make(map[string]interface{})
			c.MyMap[mod]["CreateTS"] = time.Now()
			c.MyMap[mod]["Status"] = "InProgress"
			c.MyMap[mod]["RETRYCOUNT"] = 0
		}

	case EGRESS_SUCCESS_RECORD:
		c.MyMap[mod]["EndTS"] = time.Now()
		c.MyMap[mod]["Status"] = "SUCCESS"
	case EGRESS_RETRY_RECORD:

		if _, ok := c.MyMap[mod]; !ok {
			fmt.Println("Something terible worng : record not created during ingress", mod)
			return
		}
		c.MyMap[mod]["EndTS"] = time.Now()
		c.MyMap[mod]["Status"] = "RETRYING"
		valueInt := c.MyMap[mod]["RETRYCOUNT"].(int)
		c.MyMap[mod]["RETRYCOUNT"] = valueInt + 1

		if valueInt > 5 {
			fmt.Println("This record had tried maximum : Mark it failure", mod)
		}
	}
	c.DumpMap()
}

func DumpValue(value interface{}) {
	switch value.(type) {
	case int:
		fmt.Println("Int:", value.(int))
	case time.Time:
		fmt.Println("Int:", value.(time.Time))
	case string:
		fmt.Println("Int:", value.(string))
	}
}
func (c *MySsh) DumpMap() {
	for key, value := range c.MyMap {
		fmt.Printf(key + ":")
		for k, v := range value {
			fmt.Printf(k)
			DumpValue(v)
		}
	}
}

func (c *MySsh) SendCommand(mod string, cmd string) {
	req := &Request{cmd: cmd + "\n", mod: mod, rcvCh: make(chan *Response)}
	fmt.Println("Sending Command ", cmd, time.Now())

	c.IngressInputCh <- req

}
func ReadUntil(sshOut io.Reader, expString string) string {
	buf := make([]byte, 1000)
	readStr := ""
	for {
		n, err := sshOut.Read(buf)
		readStr = readStr + string(buf[:n])
		if err == io.EOF || n == 0 {
			fmt.Println("Break,", err, n)
			break
		}

		if strings.Contains(readStr, expString) {
			break
		}

	}
	return readStr
}

func (c *MySsh) Final(rchan chan interface{}) chan interface{} {
	leftCh := make(chan interface{})
	go func() {

		for val := range rchan {
			session, _ := c.client.NewSession()
			req := val.(*Request)
			str, _ := session.Output(req.cmd)
			resp := &Response{output: string(str), cmd: req.cmd, mod: req.mod}
			leftCh <- resp
			session.Close()

		}
	}()
	return leftCh
}

func (c *MySsh) SaveContext(rchan chan interface{}) chan interface{} {
	leftCh := make(chan interface{})
	go func() {
		for val := range rchan {
			fmt.Println("Inside Save Context")
			leftCh <- val
		}
	}()
	return leftCh
}

func (c *MySsh) Audit(rchan chan interface{}) chan interface{} {
	leftCh := make(chan interface{})

	go func() {
		for val := range rchan {
			fmt.Println("Inside Audit")
			leftCh <- val
		}
	}()
	return leftCh
}
func (c *MySsh) RecordIngress(rchan chan interface{}) chan interface{} {
	leftCh := make(chan interface{})

	go func() {
		for val := range rchan {
			fmt.Println("Inside recordIngress")
			c.RecordMap(INGRESS_UPDATE_RECORD, val.(*Request).mod)
			leftCh <- val
		}
	}()
	return leftCh
}

func (c *MySsh) RecordEgress(rightCh chan interface{}) chan interface{} {
	leftCh := make(chan interface{})
	go func() {
		for val := range rightCh {
			fmt.Println("Inside Record Egress, timestamp on Map ")
			leftCh <- val
		}
	}()
	return leftCh
}

func (c *MySsh) ProcessRegExpReponse(leftCh chan interface{}, val interface{}) {
	// If the command is not success, call sendCommand Again - After 2 seconds/3seconds etc
	// Update Failure Count
	// If more than required number, then stop doing this
	status := false

	defer func() {
		if status {
			// On Success, inform the Final Go Routine
			c.RecordMap(EGRESS_SUCCESS_RECORD, val.(*Response).mod)
			leftCh <- "Success"
		} else {
			// This create a new go routine - for recycling on failure
			fmt.Println("Recycling the command again:", val.(*Response).mod)
			time.Sleep(time.Second * c.RecycleDuration)
			c.RecordMap(EGRESS_RETRY_RECORD, val.(*Response).mod)
			c.SendCommand(val.(*Response).mod, val.(*Response).cmd)
		}
	}()

	str := val.(*Response).output
	rexp, _ := regexp.Compile(`\s+(\d+)\s+`)
	slice := rexp.FindStringSubmatch(str)
	if len(slice) == 0 {
		fmt.Println("REgular expression ERROR:", str)
		return
	}
	if value, err := strconv.Atoi(slice[1]); err == nil {
		if value == 1 {
			status = true
		}
	}

	return
}
func (c *MySsh) ProcessEgress(rchan chan interface{}) chan interface{} {
	leftCh := make(chan interface{})

	go func() {
		for val := range rchan {
			fmt.Println("Check state in PRocessegress ", val.(*Response).cmd)
			go c.ProcessRegExpReponse(leftCh, val)
		}
	}()
	return leftCh
}

func (c *MySsh) FinalEgress(rchan chan interface{}) chan interface{} {
	leftCh := make(chan interface{})

	go func() {
		for val := range rchan {
			fmt.Println("Inside FinalEgress, ", val)
			leftCh <- val
		}
	}()
	return leftCh
}
