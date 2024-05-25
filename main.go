package main

import (
	"fmt"
	"math/rand"
	"os"
	"time"
)

const internalCommunicationChannel = "127.0.0.1:20202"

func main() {
	if len(os.Args) < 2 {
		fmt.Println("start <My\\Path\\To\\Config>|listen|check|end")
		os.Exit(1)
	}

	core(os.Args[1])

}

func actionLoop() string {
	time.Sleep(100 * time.Millisecond)
	return fmt.Sprintf("%d\n", rand.Intn(1000))
}
