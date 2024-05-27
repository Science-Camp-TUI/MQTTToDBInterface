package main

import (
	"fmt"
	"os"
)

const (
	Reset = "\033[0m"
	Red   = "\033[31m"
	Green = "\033[32m"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("MQTTToDB <My\\Path\\To\\Config>")
		os.Exit(1)
	}
	core(os.Args[1])

}
