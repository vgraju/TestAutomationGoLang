package Logging

import "fmt"

type Log struct {
}

func Init() *Log {
	fmt.Println("Logging enabled")
	return &Log{}
}
func (log *Log) Debug(str string) {
	fmt.Println(str)
}
