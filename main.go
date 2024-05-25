package main

import (
	"fmt"
	"os"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("MQTTToDB <My\\Path\\To\\Config>")
		os.Exit(1)
	}

	core(os.Args[1])

}
